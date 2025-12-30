package filter

type RuleType string

const (
	Basename RuleType = "basename"
	Path     RuleType = "path"
)

type Rule struct {
	Pattern string   `yaml:"pattern"`
	Type    RuleType `yaml:"type"`
}

