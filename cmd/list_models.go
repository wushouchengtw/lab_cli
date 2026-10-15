package cmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path"

	"lab_cli/service"
)

// ListModelsOptions contains the flags for the list-models command.
type ListModelsOptions struct {
	Board string
	Key   string
	Limit int
}

// parseListModelsFlags parses and validates CLI arguments for list-models.
func parseListModelsFlags(args []string) (*ListModelsOptions, error) {
	fs := flag.NewFlagSet("list-models", flag.ContinueOnError)
	opts := &ListModelsOptions{}

	fs.StringVar(&opts.Board, "b", "", "Board name (required)")
	fs.StringVar(&opts.Board, "board", "", "Board name (alias for -b)")
	fs.IntVar(&opts.Limit, "limit", 1000, "Max number of models to return (optional)")
	fs.StringVar(&opts.Key, "k", "", "Path to service account key file (optional)")
	fs.StringVar(&opts.Key, "key", "", "Path to service account key file (optional)")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s lab_cli list-models -b <board> [options]\n\n%s\n", Header("Usage:"), Header("Options:"))
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if opts.Board == "" {
		fs.Usage()
		return nil, fmt.Errorf("flag -b (board) is required")
	}

	return opts, nil
}

// RunListModels executes the list-models subcommand.
func RunListModels(args []string) error {
	opts, err := parseListModelsFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	bs, err := service.New(ctx, opts.Key)
	if err != nil {
		return fmt.Errorf("failed to initialize build service: %w", err)
	}

	models, err := bs.ListModels(ctx, opts.Board, opts.Limit)
	if err != nil {
		return fmt.Errorf("failed to list models for board %q: %w", opts.Board, err)
	}

	if len(models) == 0 {
		slog.Debug("No models found for board", "Board", opts.Board)
		return nil
	}

	if opts.Limit > 0 && len(models) > opts.Limit {
		models = models[:opts.Limit]
	}

	slog.Debug("Found model(s) for board", "Board", opts.Board, "Count", len(models))
	for _, m := range models {
		modelName := path.Base(m)
		slog.Debug("Model", "ModelName", modelName)
	}

	return nil
}
