package manifest

import "fmt"

func (v InvariantViolation) Error() string {
	if v.Path != "" {
		return fmt.Sprintf(
			"%s: %s (%s)",
			v.Path,
			v.Invariant,
			v.Capability,
		)
	}

	return fmt.Sprintf(
		"%s (%s)",
		v.Invariant,
		v.Capability,
	)
}


func (m *Manifest) validateCapabilities() error {
	declared := m.Manifest.Capabilities.Rollup.Declared()
	violations := m.collectRollupCapabilityViolations()

	for capName, enabled := range declared {
		if !enabled {
			continue
		}

		invariants, ok := rollupCapabilityInvariants[capName]
		if !ok {
			// Unknown or future capability — ignore per spec
			continue
		}

		if enabled && len(invariants) == 0 {
			return fmt.Errorf("declared capability %s has no invariants", capName)
		}

		for _, n := range m.Nodes {
			if !n.IsDir || n.Rollup == nil {
				continue
			}

			for _, inv := range invariants {
				if err := inv.Validate(n); err != nil {
					return fmt.Errorf(
						"%s: capability %s invariant %s violated: %w",
						n.Path,
						capName,
						inv.Name,
						err,
					)
				}
			}
		}
	}

	for _, v := range violations {
		if v.IsFatal() {
			return v
		}
	}

	return nil
}

