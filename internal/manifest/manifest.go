package manifest

import (
	"time"
)

type Manifest struct {
    Root      string      `json:"root"`
    Generated time.Time   `json:"generated_at"`
    Nodes     []*Node     `json:"nodes"`
    Skipped   []SkippedEntry `json:"skipped,omitempty"`
}


type SkippedEntry struct {
    Path   string `json:"path"`
    IsDir  bool   `json:"is_dir"`
    Reason string `json:"reason"`
    Rule   string `json:"rule,omitempty"`
}

func (m *Manifest) PrettySkipped() []string {
	output := make([]string, len(m.Skipped))

	for i,s := range m.Skipped {
		output[i] = s.Path
	}

	return output
}
