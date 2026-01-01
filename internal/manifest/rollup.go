package manifest

import (
	"fmt"
	"path/filepath"
	"sort"
)

type RollupOptions struct {
	EnableDirCounts   bool
	EnableSizeBytes   bool
	EnableFileTypes   bool
	EnablePercentiles bool

	// TODO: Future work
	EnableDepthStats  bool
}

type Rollup struct {
	TotalFiles   int            `json:"total_files"`
	TotalDescendantDirs   int   `json:"total_descendant_dirs"`
	Extensions   map[string]int `json:"extensions,omitempty"`

	// Size statistics in bytes
	Size struct {
		Total  int64 `json:"total"`
		Min    int64 `json:"min,omitempty"`
		Max    int64 `json:"max,omitempty"`
		Mean   int64 `json:"mean,omitempty"`
		Median int64 `json:"median,omitempty"`
		
		// p50, 90, 99
		Percentiles *Percentiles `json:"percentiles,omitempty"`
		Buckets     *SizeBuckets `json:"buckets,omitempty"`
	} `json:"size"`

	LastModified int64 `json:"last_modified"`
}

type RollupCapabilities struct {
	// Size-related
	SizeStats       bool `json:"size_stats"`
	SizePercentiles bool `json:"size_percentiles"`
	SizeBuckets     bool `json:"size_buckets"`

	// Time-related
	ActivitySpan    bool `json:"activity_span"`

	// Structure-related
	DirCounts       bool `json:"dir_counts"`
	DepthStats      bool `json:"depth_stats"`
	DepthMetrics    bool `json:"depth_metrics"`

	// Content-related
	ExtensionCounts bool `json:"extension_counts"`
	FileTypes       bool `json:"file_types"`
}

func (rc RollupCapabilities) Declared() map[string]bool {
	return map[string]bool{
		"size_stats":        rc.SizeStats,
		"size_buckets":      rc.SizeBuckets,
		"activity_span":     rc.ActivitySpan,
		// intentionally omit non-spec capabilities
	}
}

type SizeBuckets struct {
	Lt1KB     int `json:"lt_1kb"`
	KbTo1MB   int `json:"kb_to_1mb"`
	MbTo10MB  int `json:"mb_to_10mb"`
	Gt10MB    int `json:"gt_10mb"`
}

type Percentiles struct {
	P50 int64 `json:"p50,omitempty"`
	P90 int64 `json:"p90,omitempty"`
	P99 int64 `json:"p99,omitempty"`
}

type ValidateOptions struct {
	Strict   bool
}

func (m *Manifest) BuildRollups(opts RollupOptions) error {
	if len(m.Nodes) == 0 {
		return nil
	}

	/*
	// Unused for now - might come back
	// 1. Index nodes by path
	nodesByPath := make(map[string]*Node, len(m.Nodes))
	for _, n := range m.Nodes {
		nodesByPath[n.Path] = n
	}
	*/

	// 2. Build parent → children map
	children := make(map[string][]*Node)

	for _, n := range m.Nodes {
	    parent := filepath.Dir(n.Path)
	    children[parent] = append(children[parent], n)
	}

	for _, n := range m.Nodes {
	    if !n.IsDir {
	        continue
	    }

	    count := 0
	    for _, c := range children[n.Path] {
	        if c.IsDir {
	            count++
	        }
	    }
	    n.DirectSubdirCount = count
	}


	// 3. Sort directories deepest-first
	dirs := make([]*Node, 0)
	for _, n := range m.Nodes {
		if n.IsDir {
			dirs = append(dirs, n)
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return depth(dirs[i].Path) > depth(dirs[j].Path)
	})

	// 4. Build rollups bottom-up
	for _, dir := range dirs {
		r := &Rollup{}

		var (
			sizeSamples []int64
			lastMod     int64
		)

		for _, child := range children[dir.Path] {
			if child.IsDir {
			    if opts.EnableDirCounts {
        			r.TotalDescendantDirs++ // direct child
	    		    if child.Rollup != nil {
    	       			r.TotalDescendantDirs += child.Rollup.TotalDescendantDirs
        			}
    			}
			} else {
				r.TotalFiles++

				if opts.EnableFileTypes {
					ext := filepath.Ext(child.Path)
					if ext != "" {
						if r.Extensions == nil {
							r.Extensions = make(map[string]int)
						}
						r.Extensions[ext]++
					}
				}

				if opts.EnableSizeBytes && child.SizeBytes > 0 {
					sizeSamples = append(sizeSamples, child.SizeBytes)
					r.Size.Total += child.SizeBytes
				}
			}

			if child.MtimeUnix > lastMod {
				lastMod = child.MtimeUnix
			}
		}

		// 5. Finalize size stats
		if opts.EnableSizeBytes && len(sizeSamples) > 0 {
		    sort.Slice(sizeSamples, func(i, j int) bool {
		        return sizeSamples[i] < sizeSamples[j]
		    })

		    r.Size.Min = sizeSamples[0]
		    r.Size.Max = sizeSamples[len(sizeSamples)-1]
		    r.Size.Mean = r.Size.Total / int64(len(sizeSamples))

		    if opts.EnablePercentiles {
		        r.Size.Percentiles = computePercentiles(sizeSamples)
		        r.Size.Median = r.Size.Percentiles.P50
		    } else {
		        r.Size.Median = median(sizeSamples)
		    }
		}
		r.LastModified = lastMod

		m.Manifest.Capabilities.Rollup = RollupCapabilities {
		    SizeStats:    opts.EnableSizeBytes,
		    SizeBuckets:  false, // future
		    //SizePercentiles:  opts.EnablePercentiles,
		    ActivitySpan: false, // future
		    FileTypes:    opts.EnableFileTypes,
		    DepthMetrics: opts.EnableDepthStats,
		}

		// Attach rollup
		dir.Rollup = r
	}

	return nil
}

func validateNode(n *Node, opts ValidateOptions) error {
	if n.Rollup != nil {
		if err := validateRollup(n); err != nil {
			return err
		}
	}
	return nil
}

func validateRollup(n *Node) error {
	r := n.Rollup
	if r == nil {
		return nil
	}

	if r.TotalFiles < n.FileCount {
		return fmt.Errorf("%s: total_files < file_count", n.Path)
	}

	if r.TotalDescendantDirs < n.DirectSubdirCount {
		return fmt.Errorf("%s: total_descendant_dirs < direct_subdir_count", n.Path)
	}
	if r.Size.Percentiles != nil {
	    p := r.Size.Percentiles

	    if r.Size.Median != p.P50 {
	        return fmt.Errorf("%s: median != p50", n.Path)
	    }

	    if !(r.Size.Min <= p.P50 &&
	        p.P50 <= p.P90 &&
	        p.P90 <= p.P99 &&
	        p.P99 <= r.Size.Max) {
	        return fmt.Errorf("%s: percentile ordering violated", n.Path)
	    }
	}

	return nil
}

