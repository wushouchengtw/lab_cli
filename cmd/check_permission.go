package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"lab_cli/service"
	"os"
	"strings"
)

// CheckPermissionOptions contains the flags for the check-permission command.
type CheckPermissionOptions struct {
	Board  string
	Models string
	Key    string
	Limit  int
	JSON   bool
}

// parseCheckPermissionFlags parses and validates CLI arguments for check-permission.
func parseCheckPermissionFlags(args []string) (*CheckPermissionOptions, error) {
	fs := flag.NewFlagSet("check-permission", flag.ContinueOnError)
	opts := &CheckPermissionOptions{}

	fs.StringVar(&opts.Board, "b", "", "Board name (required)")
	fs.StringVar(&opts.Board, "board", "", "Board name (alias for -b)")

	fs.StringVar(&opts.Models, "ms", "", "Models name (required)")
	fs.StringVar(&opts.Models, "models", "", "Models name (alias for -ms)")

	fs.IntVar(&opts.Limit, "limit", 1000, "Max number of models to return (optional)")
	fs.StringVar(&opts.Key, "k", "", "Path to service account key file (optional)")
	fs.StringVar(&opts.Key, "key", "", "Path to service account key file (optional)")

	fs.BoolVar(&opts.JSON, "json", false, "Output in JSON format")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s lab_cli check-permission -b <board>\n\n%s\n", Header("Usage:"), Header("Options:"))
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

// RunCheckPermission executes the check-permission subcommand.
func RunCheckPermission(args []string) error {
	opts, err := parseCheckPermissionFlags(args)
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

	models := strings.Split(opts.Models, ",")
	for i, model := range models {
		models[i] = strings.TrimSpace(model)
	}

	modelsWithPermission, err := bs.CheckPermission(ctx, opts.Board, models, opts.Limit)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}

	if opts.JSON {
		return json.NewEncoder(os.Stdout).Encode(modelsWithPermission)
	}

	for _, model := range modelsWithPermission {
		status := "False"
		if model.HasPermission {
			status = "True"
		}
		fmt.Fprintf(os.Stdout, "%s: %s\n", model.Model, status)
	}
	return nil
}
