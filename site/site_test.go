package site_test

import (
	"testing"

	"lab_cli/site"
)

func TestGetServiceAccountPath(t *testing.T) {
	t.Run("returns injected key when provided", func(t *testing.T) {
		t.Setenv(site.ServiceAccountKeyPathEnv, "/env/path/key.json")
		got := site.GetServiceAccountPath("/injected/path/key.json")
		want := "/injected/path/key.json"
		if got != want {
			t.Errorf("GetServiceAccountPath() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to env var when no injected key", func(t *testing.T) {
		t.Setenv(site.ServiceAccountKeyPathEnv, "/env/path/key.json")
		got := site.GetServiceAccountPath()
		want := "/env/path/key.json"
		if got != want {
			t.Errorf("GetServiceAccountPath() = %q, want %q", got, want)
		}
	})

	t.Run("falls back to env var when injected key is empty string", func(t *testing.T) {
		t.Setenv(site.ServiceAccountKeyPathEnv, "/env/path/key.json")
		got := site.GetServiceAccountPath("")
		want := "/env/path/key.json"
		if got != want {
			t.Errorf("GetServiceAccountPath() = %q, want %q", got, want)
		}
	})

	t.Run("returns empty string when neither is provided", func(t *testing.T) {
		t.Setenv(site.ServiceAccountKeyPathEnv, "")
		got := site.GetServiceAccountPath()
		want := ""
		if got != want {
			t.Errorf("GetServiceAccountPath() = %q, want %q", got, want)
		}
	})
}
