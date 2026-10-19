package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

func RunBuild(args []string) error {
	if len(args) == 0 {
		PrintBuildUsage()
		return fmt.Errorf("no build subcommand specified")
	}
	subcommand := args[0]
	subArgs := args[1:]
	switch subcommand {
	case "get":
		return RunBuildGet(subArgs)
	case "cancel":
		return RunBuildCancel(subArgs)
	case "-h", "--help", "help":
		PrintBuildUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown build command: %q\n\n", subcommand)
		PrintBuildUsage()
		return fmt.Errorf("unknown build command: %s", subcommand)
	}
}

func PrintBuildUsage() {
	NewHelpBuilder().
		SetUsage("lab_cli build <command> [options]").
		AddCommand("get", "Get detailed build information").
		AddCommand("cancel", "Cancel a running build").
		AddOption("-h, --help", "Show help message").
		SetHintCommand("lab_cli build <command>").
		Print()
}

type BuildGetOptions struct {
	ID int64
}

func parseBuildGetFlags(args []string) (*BuildGetOptions, error) {
	fs := flag.NewFlagSet("build get", flag.ContinueOnError)
	opts := &BuildGetOptions{}
	fs.Int64Var(&opts.ID, "id", 0, "Build ID (required)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s lab_cli build get --id <ID> [options]\n\n%s\n", Header("Usage:"), Header("Options:"))
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if opts.ID == 0 {
		fs.Usage()
		return nil, fmt.Errorf("flag --id is required")
	}
	return opts, nil
}

func RunBuildGet(args []string) error {
	opts, err := parseBuildGetFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	// TODO:
	_ = opts
	return nil
}

type BuildCancelOptions struct {
	ID int64
}

func parseBuildCancelFlags(args []string) (*BuildCancelOptions, error) {
	fs := flag.NewFlagSet("build cancel", flag.ContinueOnError)
	opts := &BuildCancelOptions{}
	fs.Int64Var(&opts.ID, "id", 0, "Build ID (required)")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s lab_cli build cancel --id <ID> [options]\n\n%s\n", Header("Usage:"), Header("Options:"))
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if opts.ID == 0 {
		fs.Usage()
		return nil, fmt.Errorf("flag --id is required")
	}
	return opts, nil
}

func RunBuildCancel(args []string) error {
	opts, err := parseBuildCancelFlags(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}
	// TODO:
	_ = opts
	return nil
}
