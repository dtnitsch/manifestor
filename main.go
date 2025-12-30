package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dtnitsch/manifestor/internal/config"
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

    sc := scanner.New(scanner.Options{
        Root:              cfg.Scanner.Root,
        MaxWorkers:        cfg.Scanner.MaxWorkers,
        FollowSymlinks:    false,
        CollectInodes:     true,
        CollectTimestamps: true,
        CollectFileCounts: true,
    }, filters)

    manifest, err := sc.Scan(context.Background())
    if err != nil {
        return err
    }

	// Double check skipped things
	skippedCount := len(manifest.Skipped)
	if skippedCount > 0 {
		logger.Info("skipped directories", "count", skippedCount, "list", manifest.PrettySkipped())

		if err := scanner.AssertNoSkippedChildLeakage(manifest); err != nil {
			return fmt.Errorf("skipped children: %w", err)
		}
	}

    return output.WriteJSON(cfg.Output.File, manifest)
}

