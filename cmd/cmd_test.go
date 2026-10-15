package cmd

import (
	"errors"
	"flag"
	"testing"
)

func TestParseCheckPermissionFlags(t *testing.T) {
	t.Run("successfully parses valid required flags", func(t *testing.T) {
		args := []string{"-b", "volteer", "-ms", "model-a"}
		opts, err := parseCheckPermissionFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Board != "volteer" {
			t.Errorf("opts.Board = %q, want %q", opts.Board, "volteer")
		}
		if opts.Models != "model-a" {
			t.Errorf("opts.Models = %q, want %q", opts.Models, "model-a")
		}
		if opts.Limit != 1000 {
			t.Errorf("opts.Limit = %d, want default 1000", opts.Limit)
		}
		if opts.JSON {
			t.Errorf("opts.JSON = true, want default false")
		}
		if opts.Key != "" {
			t.Errorf("opts.Key = %q, want empty", opts.Key)
		}
	})

	t.Run("successfully parses all optional flags", func(t *testing.T) {
		args := []string{
			"--board", "brya",
			"--models", "model-b",
			"--limit", "50",
			"--key", "/path/to/key.json",
			"--json",
		}
		opts, err := parseCheckPermissionFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Board != "brya" {
			t.Errorf("opts.Board = %q, want %q", opts.Board, "brya")
		}
		if opts.Models != "model-b" {
			t.Errorf("opts.Models = %q, want %q", opts.Models, "model-b")
		}
		if opts.Limit != 50 {
			t.Errorf("opts.Limit = %d, want 50", opts.Limit)
		}
		if opts.Key != "/path/to/key.json" {
			t.Errorf("opts.Key = %q, want %q", opts.Key, "/path/to/key.json")
		}
		if !opts.JSON {
			t.Errorf("opts.JSON = false, want true")
		}
	})

	t.Run("returns error when board is missing", func(t *testing.T) {
		args := []string{"-ms", "model-a"}
		_, err := parseCheckPermissionFlags(args)

		if err == nil {
			t.Fatal("expected error when board flag is missing, got nil")
		}
	})

	t.Run("returns ErrHelp when -h is passed", func(t *testing.T) {
		args := []string{"-h"}
		_, err := parseCheckPermissionFlags(args)

		if !errors.Is(err, flag.ErrHelp) {
			t.Errorf("expected ErrHelp, got %v", err)
		}
	})
}

func TestParseListModelsFlags(t *testing.T) {
	t.Run("successfully parses minimal required flags", func(t *testing.T) {
		args := []string{"-b", "volteer"}
		opts, err := parseListModelsFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Board != "volteer" {
			t.Errorf("opts.Board = %q, want %q", opts.Board, "volteer")
		}
		if opts.Limit != 1000 {
			t.Errorf("opts.Limit = %d, want default 1000", opts.Limit)
		}
		if opts.Key != "" {
			t.Errorf("opts.Key = %q, want empty", opts.Key)
		}
	})

	t.Run("successfully parses all optional flags", func(t *testing.T) {
		args := []string{
			"--board", "brya",
			"--limit", "25",
			"-k", "/path/to/key.json",
		}
		opts, err := parseListModelsFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Board != "brya" {
			t.Errorf("opts.Board = %q, want %q", opts.Board, "brya")
		}
		if opts.Limit != 25 {
			t.Errorf("opts.Limit = %d, want 25", opts.Limit)
		}
		if opts.Key != "/path/to/key.json" {
			t.Errorf("opts.Key = %q, want %q", opts.Key, "/path/to/key.json")
		}
	})

	t.Run("returns error when board is missing", func(t *testing.T) {
		args := []string{"--limit", "10"}
		_, err := parseListModelsFlags(args)

		if err == nil {
			t.Fatal("expected error when board flag is missing, got nil")
		}
	})

	t.Run("returns ErrHelp when -h is passed", func(t *testing.T) {
		args := []string{"-h"}
		_, err := parseListModelsFlags(args)

		if !errors.Is(err, flag.ErrHelp) {
			t.Errorf("expected ErrHelp, got %v", err)
		}
	})
}

func TestExecute(t *testing.T) {
	t.Run("returns error when no arguments specified", func(t *testing.T) {
		err := Execute([]string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error when unknown command specified", func(t *testing.T) {
		err := Execute([]string{"invalid-subcommand"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("handles help flags gracefully", func(t *testing.T) {
		helpArgs := [][]string{
			{"-h"},
			{"--help"},
			{"help"},
			{"-v", "-h"},
			{"--verbose", "help"},
		}

		for _, args := range helpArgs {
			if err := Execute(args); err != nil {
				t.Errorf("Execute(%v) returned unexpected error: %v", args, err)
			}
		}
	})
}

func TestHeader(t *testing.T) {
	got := Header("Title:")
	want := "\033[1;33mTitle:\033[0m"
	if got != want {
		t.Errorf("Header() = %q, want %q", got, want)
	}
}

func TestSetupLogger(t *testing.T) {
	// Ensure setupLogger handles both verbose true and false without panic
	setupLogger(false)
	setupLogger(true)
}
