package manifest

import "log/slog"

type ViolationSummary struct {
	Total    int
	Errors   int
	Warnings int

	ByCapability map[string]int
	ByInvariant  map[string]int
}

func (s ViolationSummary) HasErrors() bool {
	return s.Errors > 0
}

func SummarizeViolations(violations []InvariantViolation) ViolationSummary {
	s := ViolationSummary{
		Total:        len(violations),
		ByCapability: make(map[string]int),
		ByInvariant:  make(map[string]int),
	}

	for _, v := range violations {
		if v.Severity.IsFatal() {
			s.Errors++
		} else {
			s.Warnings++
		}

		if v.Capability != "" {
			s.ByCapability[v.Capability]++
		}

		if v.Invariant != "" {
			s.ByInvariant[v.Invariant]++
		}
	}

	return s
}

func LogViolationSummary(logger *slog.Logger, s ViolationSummary) {
	logger.Info(
		"validation summary",
		slog.Int("total", s.Total),
		slog.Int("errors", s.Errors),
		slog.Int("warnings", s.Warnings),
		slog.Any("by_capability", s.ByCapability),
	)
}

