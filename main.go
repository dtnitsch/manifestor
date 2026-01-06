package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dtnitsch/manifestor/internal/config"
	"github.com/dtnitsch/manifestor/internal/manifest"
	"github.com/dtnitsch/manifestor/internal/output"
	"github.com/dtnitsch/manifestor/internal/scanner"
	"github.com/urfave/cli/v2"
)

const version = "0.3.0"

func main() {
	app := &cli.App{
		Name:    "manifestor",
		Usage:   "Generate LLM-friendly filesystem manifests",
		Version: version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "root",
				Aliases: []string{"r"},
				Usage:   "Root directory to scan (overrides config)",
			},
			&cli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "Output format: json or yaml (overrides config)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path (overrides config)",
			},
			&cli.StringFlag{
				Name:  "config",
				Usage: "Config file path",
				Value: "manifestor-config.yaml",
			},
		},
		Action: func(c *cli.Context) error {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			// Load config
			configPath := c.String("config")
			cfg, err := config.Load(logger, configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Override config with CLI flags
			if c.IsSet("root") {
				cfg.Scanner.Root = c.String("root")
			}
			if c.IsSet("format") {
				cfg.Output.Format = c.String("format")
			}
			if c.IsSet("output") {
				cfg.Output.File = c.String("output")
			}

			if err := run(logger, cfg); err != nil {
				return err
			}

			logger.Info("DONE - Success!")
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cfg *config.Config) error {
    filters := scanner.FilterSet{
        Block: cfg.Filters.Block,
        Allow: cfg.Filters.Allow,
    }

	opts := scanner.Options{
        Root:              cfg.Scanner.Root,
        MaxWorkers:        cfg.Scanner.MaxWorkers,
        FollowSymlinks:    false,
        CollectInodes:     true,
        CollectTimestamps: true,
        CollectFileCounts: true,
	}

    sc := scanner.New(opts, filters)

    m, err := sc.Scan(context.Background())
    if err != nil {
        return err
    }

	// Double check skipped things
	skippedCount := len(m.Skipped)
	if skippedCount > 0 {
		logger.Info("skipped directories", "count", skippedCount, "list", m.PrettySkipped())

		if err := scanner.AssertNoSkippedChildLeakage(m); err != nil {
			return fmt.Errorf("skipped children: %w", err)
		}
	}

	if cfg.Rollup.Enable {
		err := m.BuildRollups(manifest.RollupOptions{
			EnableDirCounts: cfg.Rollup.EnableDirCounts,
			EnableSizeBytes: cfg.Rollup.EnableSizeBytes,
			EnableFileTypes: cfg.Rollup.EnableFileTypes,
			EnableDepthStats: cfg.Rollup.EnableDepthStats,
			EnablePercentiles: cfg.Rollup.EnablePercentiles,
		})
		if err != nil {
			return fmt.Errorf("rollups: %w", err)
		}
		if cfg.Validate.Enable {
            violations, err := m.Validate(manifest.ValidateOptions{
                Strict: true,
            })
			// Violations are returned for reporting; fatal capability violations surface via err.
            if err != nil {
                return fmt.Errorf("validation failed: %w", err)
            }

            for _, v := range violations {
                if v.Severity == manifest.SeverityWarning {
					manifest.LogViolation(logger, v)
                }
            }

			summary := manifest.SummarizeViolations(violations)
			if summary.HasErrors() {
				manifest.LogViolationSummary(logger, summary)
				return fmt.Errorf("validation failed")
			}
		}
	}

	// Set defaults
	m.Manifest = manifest.DefaultManifestMeta()

	// Write output based on configured format
	switch cfg.Output.Format {
	case "yaml":
		return output.WriteYAML(cfg.Output.File, m)
	case "json":
		return output.WriteJSON(cfg.Output.File, m)
	default:
		return fmt.Errorf("unsupported output format: %s (supported: json, yaml)", cfg.Output.Format)
	}
}

