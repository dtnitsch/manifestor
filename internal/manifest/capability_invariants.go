package manifest

import "fmt"

type RollupInvariant struct {
	Name     string
	Validate func(n *Node) error
}

var rollupCapabilityInvariants = map[string][]RollupInvariant{
	"size_stats": {
		{
			Name: "size.total",
			Validate: func(n *Node) error {
				if n.Rollup.Size.Total == 0 && n.Rollup.TotalFiles > 0 {
					return fmt.Errorf("size.total missing")
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
}

