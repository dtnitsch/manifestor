package manifest

import (
	"path/filepath"
	"sort"
)

type RollupOptions struct {
	EnableDirCounts bool
	EnableSizeBytes bool
	EnableFileTypes bool
	EnableDepthStats bool
}

type Rollup struct {
	TotalFiles   int            `json:"total_files"`
	DescendantDirs    int            `json:"descendant_dirs"`
	Extensions   map[string]int `json:"extensions,omitempty"`

	Size struct {
		Total  int64 `json:"total"`
		Min    int64 `json:"min,omitempty"`
		Max    int64 `json:"max,omitempty"`
		Mean   int64 `json:"mean,omitempty"`
		Median int64 `json:"median,omitempty"`
	} `json:"size"`

	LastModified int64 `json:"last_modified"`
}

type rollupState struct {
	files      int
	dirs       int
	bytes      int64
	extensions map[string]int
	depth      int
	sizes      []int64
	lastMod    int64
}


func (m *Manifest) BuildRollups(opts RollupOptions) error {
	if len(m.Nodes) == 0 {
		return nil
	}

	// 1. Index nodes by path
	nodesByPath := make(map[string]*Node, len(m.Nodes))
	for _, n := range m.Nodes {
		nodesByPath[n.Path] = n
	}

	// 2. Build parent → children map
	children := make(map[string][]*Node)
	for _, n := range m.Nodes {
		parent := filepath.Dir(n.Path)
		if parent == "." && n.Path == "." {
			continue
		}
		children[parent] = append(children[parent], n)
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
					r.DescendantDirs++
				}
				if child.Rollup != nil {
					r.TotalFiles += child.Rollup.TotalFiles
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

			if child.ModTimeUnix > lastMod {
				lastMod = child.ModTimeUnix
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
			r.Size.Median = median(sizeSamples)
		}

		r.LastModified = lastMod

		// Attach rollup
		dir.Rollup = r
	}

	return nil
}
