package manifest

import "fmt"

func (m *Manifest) validateCapabilities() error {
	caps := m.Manifest.Capabilities.Rollup

	for _, n := range m.Nodes {
		if !n.IsDir || n.Rollup == nil {
			continue
		}

		if err := applyRollupCapability("size_stats", caps.SizeStats, n); err != nil {
			return err
		}
		if err := applyRollupCapability("size_percentiles", caps.SizePercentiles, n); err != nil {
			return err
		}
		if err := applyRollupCapability("extension_counts", caps.ExtensionCounts, n); err != nil {
			return err
		}
	}

	return nil
}

func applyRollupCapability(name string, enabled bool, n *Node) error {
	if !enabled {
		return nil
	}

	invariants := rollupCapabilityInvariants[name]
	for _, inv := range invariants {
		if err := inv.Validate(n); err != nil {
			return fmt.Errorf(
				"%s: capability %s invariant %s failed: %w",
				n.Path,
				name,
				inv.Name,
				err,
			)
		}
	}

	return nil
}

