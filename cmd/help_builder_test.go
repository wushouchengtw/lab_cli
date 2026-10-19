package cmd

import (
	"strings"
	"testing"
)

func TestHelpBuilder(t *testing.T) {
	t.Run("builds formatted help message with aligned columns", func(t *testing.T) {
		output := NewHelpBuilder().
			SetUsage("test_cli <command> [options]").
			AddCommand("short", "Short command").
			AddCommand("longer-command", "Longer command description").
			SetOptionsTitle("Global Options:").
			AddOption("-v", "Verbose").
			AddOption("--help", "Help").
			SetHintCommand("test_cli <command>").
			String()

		// Verify sections exist
		if !strings.Contains(output, "Usage:") {
			t.Errorf("output does not contain 'Usage:' section")
		}
		if !strings.Contains(output, "Available Commands:") {
			t.Errorf("output does not contain 'Available Commands:' section")
		}
		if !strings.Contains(output, "Global Options:") {
			t.Errorf("output does not contain 'Global Options:' section")
		}
		if !strings.Contains(output, "Use \"test_cli <command> -h\" for more information") {
			t.Errorf("output does not contain expected hint")
		}

		// Verify alignment: "short" should be padded to match "longer-command"
		if !strings.Contains(output, "  short            Short command") {
			t.Errorf("output did not pad commands properly:\n%s", output)
		}
	})

	t.Run("handles minimal options without panic", func(t *testing.T) {
		output := NewHelpBuilder().
			SetUsage("minimal").
			String()

		if !strings.Contains(output, "minimal") {
			t.Errorf("expected usage in output, got: %s", output)
		}
	})

	t.Run("supports raw hint", func(t *testing.T) {
		output := NewHelpBuilder().
			SetHint("Custom hint here").
			String()

		if !strings.Contains(output, "Custom hint here") {
			t.Errorf("expected custom hint, got: %s", output)
		}
	})
}
