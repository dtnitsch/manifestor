package manifest

import "fmt"

type RollupInvariant struct {
    Name        string
    Description string
    Severity    Severity // Error | Warning (future)
    Validate    func(*Node) error
}

var rollupCapabilityInvariants = map[string][]RollupInvariant{
    "activity_span": {
        {
            Name: "activity.last_modified.present",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.LastModified == 0 {
                    return fmt.Errorf("last_modified missing")
                }
                return nil
            },
        },
        {
            Name: "activity.last_modified.bounds",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.LastModified < n.MtimeUnix {
                    return fmt.Errorf(
                        "last_modified (%d) < node.mtime_unix (%d)",
                        n.Rollup.LastModified,
                        n.MtimeUnix,
                    )
                }
                return nil
            },
        },
    },
	"dir_counts": {
		{
			Name: "dir_counts",
			Validate: func(n *Node) error {
				if n.Rollup.TotalDescendantDirs < n.DirectSubdirCount {
					return fmt.Errorf(
						"total_descendant_dirs (%d) < direct_subdir_count (%d)",
						n.Rollup.TotalDescendantDirs,
						n.DirectSubdirCount,
					)
				}
				return nil
			},
		},
	},
    "size_stats": {
        {
            Name: "size.total.present",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.TotalFiles > 0 && n.Rollup.Size.Total == 0 {
                    return fmt.Errorf("size.total missing or zero with nonzero files")
                }
                return nil
            },
        },
        {
            Name: "size.min.present",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.TotalFiles > 0 && n.Rollup.Size.Min == 0 {
                    return fmt.Errorf("size.min missing")
                }
                return nil
            },
        },
        {
            Name: "size.max.present",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.TotalFiles > 0 && n.Rollup.Size.Max == 0 {
                    return fmt.Errorf("size.max missing")
                }
                return nil
            },
        },
        {
            Name: "size.mean.present",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.TotalFiles > 0 && n.Rollup.Size.Mean == 0 {
                    return fmt.Errorf("size.mean missing")
                }
                return nil
            },
        },
        {
            Name: "size.median.present",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.TotalFiles > 0 && n.Rollup.Size.Median == 0 {
                    return fmt.Errorf("size.median missing")
                }
                return nil
            },
        },
        {
            Name: "size.ordering",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                if n.Rollup.TotalFiles == 0 {
                    return nil
                }

                min := n.Rollup.Size.Min
                median := n.Rollup.Size.Median
                max := n.Rollup.Size.Max

                if min > median || median > max {
                    return fmt.Errorf(
                        "invalid ordering: min=%d median=%d max=%d",
                        min, median, max,
                    )
                }
                return nil
            },
        },
	    {
    	    Name: "size_stats.median.matches_p50",
        	Validate: func(n *Node) error {
            	p := n.Rollup.Size.Percentiles
            	if p == nil {
                	return nil
            	}
        	    if n.Rollup.Size.Median != 0 && p.P50 != 0 && n.Rollup.Size.Median != p.P50 {
	                return fmt.Errorf("median != p50")
            	}
        	    return nil
    	    },
	    },
    },
	"size_percentiles": {
		{
			Name: "size.percentiles",
			Validate: func(n *Node) error {
				if n.Rollup.Size.Percentiles == nil {
					return fmt.Errorf("percentiles missing")
				}
				return nil
			},
		},
	    {
            Name: "size_percentiles.ordering",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                p := n.Rollup.Size.Percentiles
                if p == nil {
                    return nil
                }
                if !(p.P50 <= p.P90 && p.P90 <= p.P99) {
                    return fmt.Errorf("p50 <= p90 <= p99 violated")
                }
                return nil
            },
        },
        {
            Name: "size_percentiles.within_min_max",
            Severity: SeverityError,
            Validate: func(n *Node) error {
                s := n.Rollup.Size
                p := s.Percentiles
                if p == nil || s.Min == 0 || s.Max == 0 {
                    return nil
                }
                if p.P50 < s.Min || p.P99 > s.Max {
                    return fmt.Errorf("percentiles outside min/max")
                }
                return nil
            },
        },
        {
            Name:        "size.percentiles.missing",
            Description: "size.percentiles missing; percentile-based reasoning unavailable",
            Severity:    SeverityWarning,
            Validate: func(n *Node) error {
                if n.Rollup.Size.Percentiles == nil {
                    return fmt.Errorf("percentiles missing")
                }
                return nil
            },
        },
	},
	"extension_counts": {
		{
			Name: "extensions",
			Validate: func(n *Node) error {
				if n.Rollup.Extensions == nil {
					return fmt.Errorf("extensions missing")
				}
				return nil
			},
		},
	},
    "size_buckets": {
      {
        Name: "size_buckets.present",
        Validate: func(n *Node) error {
          if n.Rollup.Size.Buckets == nil {
            return fmt.Errorf("size.buckets missing")
          }
          return nil
        },
      },
      {
        Name: "size_buckets.keys",
        Validate: func(n *Node) error {
          b := n.Rollup.Size.Buckets
          if b == nil {
            return nil // presence checked separately
          }

          // Struct-based buckets: keys must exist implicitly,
          // but we still guard against zero-value misuse
          if b.Lt1KB < 0 || b.KbTo1MB < 0 || b.MbTo10MB < 0 || b.Gt10MB < 0 {
            return fmt.Errorf("bucket values must be >= 0")
          }

          return nil
        },
      },
      {
        Name: "size_buckets.sum",
        Validate: func(n *Node) error {
          b := n.Rollup.Size.Buckets
          if b == nil {
            return nil
          }

          sum := b.Lt1KB + b.KbTo1MB + b.MbTo10MB + b.Gt10MB
          if sum != n.Rollup.TotalFiles {
            return fmt.Errorf(
              "bucket sum %d != total_files %d",
              sum,
              n.Rollup.TotalFiles,
            )
          }
          return nil
        },
      },
    },
}

