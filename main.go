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
)

func main() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))

	configPath := "config.yaml"
	cfg, err := config.Load(logger, configPath) 
    if err != nil {
        logger.Error("failed to load config", "error", err, "path", configPath)
        os.Exit(1)
    }

    if err := run(logger, cfg); err != nil {
        logger.Error("application error", "error", err)
        os.Exit(1)
    }

	logger.Info("DONE - Success!")
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
			if err := m.Validate(manifest.ValidateOptions{
				Strict: true,
			}); err != nil {
				return fmt.Errorf("validation failed: %w", err)
			}
		}
	}

    return output.WriteJSON(cfg.Output.File, m)
}

