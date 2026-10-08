package cmd

import (
	"fmt"
	"os"
	"time"
)

// DefaultTimeout is the default timeout duration for network and RPC calls.
const DefaultTimeout = 20 * time.Second

// Execute routes subcommands and handles global flags like -h/--help.
func Execute(args []string) error {
	if len(args) == 0 {
		PrintUsage()
		return fmt.Errorf("no command specified")
	}

	command := args[0]
	subArgs := args[1:]

	switch command {
	case "list-models":
		return RunListModels(subArgs)
	case "check-permission":
		return RunCheckPermission(subArgs)
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
	fmt.Printf(`%s
  lab_cli <command> [options]

%s
  list-models        List all models for a given board
  check-permission   Check access permission for a given board

Use "lab_cli <command> -h" for more information about a command.
`, Header("Usage:"), Header("Available Commands:"))
}
