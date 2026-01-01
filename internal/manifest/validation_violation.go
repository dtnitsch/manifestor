package manifest

type InvariantViolation struct {
	Path        string
	Capability  string
	Invariant   string
	Description string
	Severity    Severity
	Err         error
}

func (v InvariantViolation) IsFatal() bool {
	return v.Severity.IsFatal()
}

// Validate evaluates all declared capability invariants.
//
// It returns:
//   - all invariant violations (fatal + warnings)
//   - a non-nil error if at least one fatal violation occurred
//
// Callers may:
//   - inspect violations for reporting or diagnostics
//   - treat error as "manifest is invalid"
func (m *Manifest) Validate(opts ValidateOptions) ([]InvariantViolation, error) {
	var allViolations []InvariantViolation

	// --- Capability-driven invariant validation ---
	capViolations, err := m.ValidateCapabilities()
	allViolations = append(allViolations, capViolations...)
	if err != nil {
		return allViolations, err
	}

	// --- Node-level validation (existing logic) ---
	for _, n := range m.Nodes {
		if err := validateNode(n, opts); err != nil {
			return allViolations, err
		}
	}

	return allViolations, nil
}

// Fatal violations are returned as the error value.
// Warnings are included in the violations slice only.
func (m *Manifest) ValidateCapabilities() ([]InvariantViolation, error) {
	violations := m.collectRollupCapabilityViolations()

	for _, v := range violations {
		if v.Severity.IsFatal() {
			return violations, v
		}
	}

	return violations, nil
}

func (m *Manifest) collectRollupCapabilityViolations() []InvariantViolation {
	var violations []InvariantViolation

	declared := m.Manifest.Capabilities.Rollup.Declared()

	for capName, enabled := range declared {
		if !enabled {
			continue
		}

		invariants, ok := rollupCapabilityInvariants[capName]
		if !ok {
			// Unknown / future capability — ignore per spec
			continue
		}

		// A declared capability with no invariants is invalid.
		// Capabilities are opt-in guarantees; zero invariants means no guarantee.
		if len(invariants) == 0 {
			violations = append(violations, InvariantViolation{
				Capability:  capName,
				Invariant:   "capability.has_invariants",
				Description: "declared capability has no invariants",
				Severity:    SeverityError,
				Err:         nil,
			})
			continue
		}

		for _, n := range m.Nodes {
			if !n.IsDir || n.Rollup == nil {
				continue
			}

			for _, inv := range invariants {
				if err := inv.Validate(n); err != nil {
					violations = append(violations, InvariantViolation{
						Path:        n.Path,
						Capability:  capName,
						Invariant:   inv.Name,
						Description: inv.Description,
						Severity:    inv.Severity,
						Err:         err,
					})
				}
			}
		}
	}

	return violations
}

