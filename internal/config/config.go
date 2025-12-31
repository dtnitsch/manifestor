package config

import (
	"log/slog"
	"os"

	yaml "gopkg.in/yaml.v3"

	"github.com/dtnitsch/manifestor/internal/filter"
)

type Config struct {
	Scanner ScannerConfig `yaml:"scanner"`
	Filters Filters       `yaml:"filters"`
	Output  Output        `yaml:"output"`
	
	// Directory Stats
	Rollup RollupConfig `yaml:"rollup"`

	// Validate
	Validate ValidateConfig `yaml:"validate"`
}

type ScannerConfig struct {
	Root           string `yaml:"root"`
	MaxWorkers     int    `yaml:"max_workers"`
	FollowSymlinks bool   `yaml:"follow_symlinks"`

	CollectInodes     bool `yaml:"collect_inodes"`
	CollectTimestamps bool `yaml:"collect_timestamps"`
	CollectFileCounts bool `yaml:"collect_file_counts"`
}

type RollupConfig struct {
	Enable   bool `yaml:"enable"`  

	EnableDirCounts bool `yaml:"enable_dir_counts"`
	EnableSizeBytes bool `yaml:"enable_size_bytes"`
	EnableFileTypes bool `yaml:"enable_file_types"`
	EnableDepthStats bool `yaml:"enable_depth_stats"`
	EnablePercentiles bool `yaml:"enable_percentiles"`
}

type ValidateConfig struct {
	Enable bool `yaml:"enable"`  
}

type Filters struct {
	Block []filter.Rule `yaml:"block"`
	Allow []filter.Rule `yaml:"allow"`
}

type Output struct {
	Format string `yaml:"format"` // json (v0.1)
	File   string `yaml:"file"`
}

func Load(log *slog.Logger, filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	applyDefaults(&cfg)
	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.Scanner.Root == "" {
		cfg.Scanner.Root = "."
	}
	if cfg.Scanner.MaxWorkers == 0 {
		cfg.Scanner.MaxWorkers = 8
	}
	if cfg.Output.Format == "" {
		cfg.Output.Format = "json"
	}
	if cfg.Output.File == "" {
		cfg.Output.File = "manifest.json"
	}
}

