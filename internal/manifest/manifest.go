package manifest

import (
	"time"
)

type Manifest struct {
	Manifest  ManifestMeta   `json:"manifest" yaml:"manifest"`
    Root      string         `json:"root" yaml:"root"`
    Generated time.Time      `json:"generated_at" yaml:"generated_at"`
    Nodes     []*Node        `json:"nodes" yaml:"nodes"`
    Skipped   []SkippedEntry `json:"skipped,omitempty" yaml:"skipped,omitempty"`
}

type SkippedEntry struct {
    Path   string `json:"path" yaml:"path"`
    IsDir  bool   `json:"is_dir,omitempty" yaml:"is_dir,omitempty"`
    Reason string `json:"reason" yaml:"reason"`
    Rule   string `json:"rule,omitempty" yaml:"rule,omitempty"`
}

type ManifestMeta struct {
	Version   string        `json:"version" yaml:"version"`
	Generator GeneratorMeta `json:"generator" yaml:"generator"`
	Schema    SchemaMeta    `json:"schema" yaml:"schema"`
	Capabilities Capabilities `json:"capabilities" yaml:"capabilities"`
}

type GeneratorMeta struct {
	Name      string `json:"name" yaml:"name"`
	Version   string `json:"version" yaml:"version"`
	Commit    string `json:"commit,omitempty" yaml:"commit,omitempty"`
	BuildTime string `json:"build_time" yaml:"build_time"`
}

type SchemaMeta struct {
	Node   string `json:"node" yaml:"node"`
	Rollup string `json:"rollup" yaml:"rollup"`
}


func (m *Manifest) PrettySkipped() []string {
	output := make([]string, len(m.Skipped))

	for i,s := range m.Skipped {
		output[i] = s.Path
	}

	return output
}

