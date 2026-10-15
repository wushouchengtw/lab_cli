package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"lab_cli/service"
)

// CheckPermissionOptions contains the flags for the check-permission command.
type CheckPermissionOptions struct {
	Board  string
	Models string
	Key    string
	Limit  int
	JSON   bool
}

// Define response struct for check permission CLI.
type CheckPermissionResponse struct {
	Board         string `json:"board"`
	Model         string `json:"model"`
	HasPermission bool   `json:"has_permission"`
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

	permission, err := bs.CheckPermission(ctx, opts.Board, opts.Models, opts.Limit)
	if err != nil {
		return fmt.Errorf("failed to check permission: %w", err)
	}

	if opts.JSON {
		res := CheckPermissionResponse{
			Board:         opts.Board,
			Model:         opts.Models,
			HasPermission: permission,
		}
		return json.NewEncoder(os.Stdout).Encode(res)
	}

	if permission {
		fmt.Fprintln(os.Stdout, "True")
	} else {
		fmt.Fprintln(os.Stdout, "False")
	}
	return nil
}
