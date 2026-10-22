package cmd

import (
	"errors"
	"flag"
	"testing"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestParseBuildGetFlags(t *testing.T) {
	t.Run("successfully parses valid required flags", func(t *testing.T) {
		args := []string{"-id", "12345"}
		opts, err := parseBuildGetFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.ID != 12345 {
			t.Errorf("opts.ID = %d, want 12345", opts.ID)
		}
		if opts.Key != "" {
			t.Errorf("opts.Key = %q, want empty", opts.Key)
		}
	})

	t.Run("successfully parses all optional flags", func(t *testing.T) {
		args := []string{"--id", "67890", "--key", "/path/to/key.json"}
		opts, err := parseBuildGetFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.ID != 67890 {
			t.Errorf("opts.ID = %d, want 67890", opts.ID)
		}
		if opts.Key != "/path/to/key.json" {
			t.Errorf("opts.Key = %q, want %q", opts.Key, "/path/to/key.json")
		}
	})

	t.Run("returns error when id flag is missing", func(t *testing.T) {
		args := []string{"--key", "/path/to/key.json"}
		_, err := parseBuildGetFlags(args)

		if err == nil {
			t.Fatal("expected error when id flag is missing, got nil")
		}
	})

	t.Run("returns ErrHelp when -h is passed", func(t *testing.T) {
		args := []string{"-h"}
		_, err := parseBuildGetFlags(args)

		if !errors.Is(err, flag.ErrHelp) {
			t.Errorf("expected ErrHelp, got %v", err)
		}
	})
}

func TestParseBuildCancelFlags(t *testing.T) {
	t.Run("successfully parses valid required flags", func(t *testing.T) {
		args := []string{"-id", "99999", "-reason", "cleanup build"}
		opts, err := parseBuildCancelFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.ID != 99999 {
			t.Errorf("opts.ID = %d, want 99999", opts.ID)
		}
		if opts.Reason != "cleanup build" {
			t.Errorf("opts.Reason = %q, want %q", opts.Reason, "cleanup build")
		}
		if opts.Key != "" {
			t.Errorf("opts.Key = %q, want empty", opts.Key)
		}
	})

	t.Run("successfully parses all optional flags", func(t *testing.T) {
		args := []string{"--id", "99999", "--reason", "cleanup build", "--key", "/path/to/key.json"}
		opts, err := parseBuildCancelFlags(args)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if opts.Key != "/path/to/key.json" {
			t.Errorf("opts.Key = %q, want %q", opts.Key, "/path/to/key.json")
		}
	})

	t.Run("returns error when id flag is missing", func(t *testing.T) {
		args := []string{"-reason", "cleanup"}
		_, err := parseBuildCancelFlags(args)

		if err == nil {
			t.Fatal("expected error when id flag is missing, got nil")
		}
	})

	t.Run("returns error when reason flag is missing", func(t *testing.T) {
		args := []string{"-id", "99999"}
		_, err := parseBuildCancelFlags(args)

		if err == nil {
			t.Fatal("expected error when reason flag is missing, got nil")
		}
	})

	t.Run("returns ErrHelp when -h is passed", func(t *testing.T) {
		args := []string{"-h"}
		_, err := parseBuildCancelFlags(args)

		if !errors.Is(err, flag.ErrHelp) {
			t.Errorf("expected ErrHelp, got %v", err)
		}
	})
}

func TestRunBuild(t *testing.T) {
	t.Run("returns error when no subcommand specified", func(t *testing.T) {
		err := RunBuild([]string{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("returns error when unknown subcommand specified", func(t *testing.T) {
		err := RunBuild([]string{"invalid-subcmd"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("handles help flags gracefully", func(t *testing.T) {
		helpArgs := [][]string{
			{"-h"},
			{"--help"},
			{"help"},
		}

		for _, args := range helpArgs {
			if err := RunBuild(args); err != nil {
				t.Errorf("RunBuild(%v) returned unexpected error: %v", args, err)
			}
		}
	})

	t.Run("handles get help gracefully", func(t *testing.T) {
		if err := RunBuild([]string{"get", "-h"}); err != nil {
			t.Errorf("RunBuild(get -h) returned unexpected error: %v", err)
		}
	})

	t.Run("handles cancel help gracefully", func(t *testing.T) {
		if err := RunBuild([]string{"cancel", "-h"}); err != nil {
			t.Errorf("RunBuild(cancel -h) returned unexpected error: %v", err)
		}
	})

	t.Run("returns error on cancel when required flags missing", func(t *testing.T) {
		if err := RunBuild([]string{"cancel", "--id", "123"}); err == nil {
			t.Fatal("expected error for missing --reason flag, got nil")
		}
	})
}

func TestToTime(t *testing.T) {
	t.Run("returns nil when proto timestamp is nil", func(t *testing.T) {
		if got := toTime(nil); got != nil {
			t.Errorf("toTime(nil) = %v, want nil", got)
		}
	})

	t.Run("returns valid Time when proto timestamp is provided", func(t *testing.T) {
		got := toTime(&timestamppb.Timestamp{
			Seconds: 0,
			Nanos:   0,
		})
		if got == nil {
			t.Fatal("expected non-nil Time, got nil")
		}
		if got.Seconds != 0 || got.TaipeiTime != "1970-01-01 08:00:00" {
			t.Errorf("unexpected Time result: %+v", got)
		}
	})
}
