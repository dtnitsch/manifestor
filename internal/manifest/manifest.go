package manifest

import (
	"time"
)

type Manifest struct {
	Manifest  ManifestMeta   `json:"manifest"`
    Root      string         `json:"root"`
    Generated time.Time      `json:"generated_at"`
    Nodes     []*Node        `json:"nodes"`
    Skipped   []SkippedEntry `json:"skipped,omitempty"`
}

type SkippedEntry struct {
    Path   string `json:"path"`
    IsDir  bool   `json:"is_dir"`
    Reason string `json:"reason"`
    Rule   string `json:"rule,omitempty"`
}

type ManifestMeta struct {
	Version   string        `json:"version"`
	Generator GeneratorMeta `json:"generator"`
	Schema    SchemaMeta    `json:"schema"`
}

type GeneratorMeta struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	BuildTime string `json:"build_time"`
}

type SchemaMeta struct {
	Node   string `json:"node"`
	Rollup string `json:"rollup"`
}

func (m *Manifest) PrettySkipped() []string {
	output := make([]string, len(m.Skipped))

	for i,s := range m.Skipped {
		output[i] = s.Path
	}

	return output
}

