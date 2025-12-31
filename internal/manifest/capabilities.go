package manifest

import (
	"fmt"
)

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

type Capabilities struct {
    Rollup RollupCapabilities `json:"rollup"`
}


func (m *Manifest) validateCapabilities() error {
	c := m.Manifest.Capabilities.Rollup

	for _, n := range m.Nodes {
		if !n.IsDir || n.Rollup == nil {
			continue
		}

		r := n.Rollup

		// ---- Size stats ----
		if c.SizeStats {
			if r.Size.Total == 0 && r.TotalFiles > 0 {
				return fmt.Errorf(
					"%s: capability size_stats claimed but size.total missing",
					n.Path,
				)
			}
		}

		// ---- Percentiles ----
		if r.Size.Percentiles == nil && r.Size.Buckets == nil {
			if r.Size.Percentiles == nil {
				return fmt.Errorf(
					"%s: capability size_percentiles claimed but size.percentiles missing",
					n.Path,
				)
			}

			p := r.Size.Percentiles
			if p.P50 == 0 || p.P90 == 0 || p.P99 == 0 {
				return fmt.Errorf(
					"%s: incomplete percentiles (p50/p90/p99)",
					n.Path,
				)
			}

			if r.Size.Median != p.P50 {
				return fmt.Errorf(
					"%s: median != p50 under size_percentiles capability",
					n.Path,
				)
			}
		}

		// ---- Extension counts ----
		if c.ExtensionCounts {
			if r.Extensions == nil {
				return fmt.Errorf(
					"%s: capability extension_counts claimed but extensions missing",
					n.Path,
				)
			}
		}

		// ---- Directory counts ----
		if c.DirCounts {
			if r.TotalDescendantDirs < n.DirectSubdirCount {
				return fmt.Errorf(
					"%s: total_descendant_dirs < direct_subdir_count",
					n.Path,
				)
			}
		}

		// ---- Depth stats (future-safe) ----
		if c.DepthStats {
			// Intentionally empty for now
			// When implemented, assert fields exist here
		}

		// ---- Size buckets ----
		if c.SizeBuckets {
			if r.Size.Buckets == nil {
				return fmt.Errorf(
					"%s: capability size_buckets claimed but buckets missing",
					n.Path,
				)
			}

			b := r.Size.Buckets
			if b.Lt1KB+b.KbTo1MB+b.MbTo10MB+b.Gt10MB != r.TotalFiles {
				return fmt.Errorf(
					"%s: size_buckets do not sum to total_files",
					n.Path,
				)
			}
		}


	}

	return nil
}

