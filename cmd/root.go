package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"time"
)

// DefaultTimeout is the default timeout duration for network and RPC calls.
const DefaultTimeout = 20 * time.Second

func setupLogger(verbose bool) {
	level := slog.LevelWarn
	if verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	slog.SetDefault(slog.New(handler))
}

// Execute routes subcommands and handles global flags like -h/--help.
func Execute(args []string) error {
	var verbose bool
	var filteredArgs []string
	for _, arg := range args {
		if arg == "-v" || arg == "--verbose" {
			verbose = true
		} else {
			filteredArgs = append(filteredArgs, arg)
		}
	}
	setupLogger(verbose)

	if len(filteredArgs) == 0 {
		PrintUsage()
		return fmt.Errorf("no command specified")
	}

	command := filteredArgs[0]
	subArgs := filteredArgs[1:]

	switch command {
	case "list-models":
		return RunListModels(subArgs)
	case "check-permission":
		return RunCheckPermission(subArgs)
	case "build":
		return RunBuild(subArgs)
	case "-h", "--help", "help":
		PrintUsage()
		return nil
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %q\n\n", command)
		PrintUsage()
		return fmt.Errorf("unknown command: %s", command)
	}
}

// Header colors section headers (e.g. Usage:, Available Commands:) in bold yellow.
func Header(text string) string {
	return "\033[1;33m" + text + "\033[0m"
}

// PrintUsage prints the CLI help message.
func PrintUsage() {
	NewHelpBuilder().
		SetUsage("lab_cli <command> [options]").
		AddCommand("list-models", "List all models for a given board").
		AddCommand("check-permission", "Check access permission for a given board").
		AddCommand("build", "Inspect or manage LUCI builds").
		SetOptionsTitle("Global Options:").
		AddOption("-v, --verbose", "Enable verbose/debug logging").
		AddOption("-h, --help", "Show help message").
		SetHint("Use \"lab_cli <command> -h\" for more information about a command.").
		Print()
}
