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

func (rc RollupCapabilities) Declared() map[string]bool {
	return map[string]bool{
		"size_stats":        rc.SizeStats,
		"size_buckets":      rc.SizeBuckets,
		"activity_span":     rc.ActivitySpan,
		// intentionally omit non-spec capabilities
	}
}

