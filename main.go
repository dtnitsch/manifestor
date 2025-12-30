package main

import (
	"context"
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

    return output.WriteJSON(cfg.Output.File, manifest)
}

