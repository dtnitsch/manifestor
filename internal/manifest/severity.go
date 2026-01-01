package manifest

import (
	"context"
	"log/slog"
)
/*

SeverityError

An invariant violation that makes the manifest semantically invalid for the declared capabilities.
Properties:

- Breaks an explicit capability guarantee
- Makes downstream reasoning unsafe
- Must fail validation in strict mode
- Indicates either a bug or corruption

Examples:
- size_stats declared but size.total missing
- size_buckets declared but buckets don’t sum
- activity_span declared but last_modified < node.mtime

These are non-negotiable.

---

SeverityWarning

An invariant violation that does not invalidate correctness, but reduces precision, confidence, or future compatibility.
Properties:

- Data is still usable
- Guarantees are weakened, not broken
- Safe to proceed with caution
- Valuable signal to humans or tooling

Examples:
- Median present but percentiles missing (if percentiles are optional)
- Extension counts missing for empty directories
- A capability declared that is deprecated but still accepted
- A capability whose invariants are partially satisfied
*/

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

func (s Severity) IsFatal() bool {
	return s == SeverityError
}

func (s Severity) Valid() bool {
	return s == SeverityError || s == SeverityWarning
}

func LogViolation(logger *slog.Logger, v InvariantViolation) {
    logger.Log(
        context.Background(),
        v.Severity.LogLevel(),
        "invariant violation",
        slog.Group(
            "violation",
            "path", v.Path,
            "capability", v.Capability,
            "invariant", v.Invariant,
            "severity", v.Severity,
            "error", v.Err,
        ),
    )
}

func (s Severity) LogLevel() slog.Level {
    switch s {
    case SeverityError:
        return slog.LevelError
    case SeverityWarning:
        return slog.LevelWarn
    default:
        return slog.LevelInfo
    }
}

/*
Not used - for later if we want it
func ParseSeverity(s string) (Severity, bool) {
	switch Severity(s) {
	case SeverityError, SeverityWarning:
		return Severity(s), true
	default:
		return "", false
	}
}
*/
