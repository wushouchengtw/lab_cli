package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"lab_cli/service"

	moblabapipb "google.golang.org/genproto/googleapis/chromeos/moblab/v1beta1"
)

// mockMoblabClient directly mocks the model slice returned from client.
type mockMoblabClient struct {
	models           []*moblabapipb.Model
	err              error
	capturedParent   string
	capturedPageSize int32
}

func (m *mockMoblabClient) ListModels(ctx context.Context, parent string, pageSize int32) ([]*moblabapipb.Model, error) {
	m.capturedParent = parent
	m.capturedPageSize = pageSize
	return m.models, m.err
}

func TestParseBuildTargetsBoardPath(t *testing.T) {
	got := service.ParseBuildTargetsBoardPath("volteer")
	want := "buildTargets/volteer"
	if got != want {
		t.Errorf("ParseBuildTargetsBoardPath() = %q, want %q", got, want)
	}
}

func TestGetModelName(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "standard resource path",
			input: "buildTargets/volteer/models/model-a",
			want:  "model-a",
		},
		{
			name:  "single model name",
			input: "model-a",
			want:  "model-a",
		},
		{
			name:  "nested path",
			input: "foo/bar/baz",
			want:  "baz",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetModelName(tt.input)
			if got != tt.want {
				t.Errorf("GetModelName(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildServiceImpl_ListModels(t *testing.T) {
	ctx := context.Background()

	t.Run("successfully lists models and parses model names", func(t *testing.T) {
		fakeModels := []*moblabapipb.Model{
			{Name: "buildTargets/volteer/models/model-a"},
			{Name: "buildTargets/volteer/models/model-b"},
		}
		mockClient := &mockMoblabClient{
			models: fakeModels,
		}

		svc := service.NewWithClient(mockClient)
		results, err := svc.ListModels(ctx, "volteer", 50)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify parameters forwarded to client
		if mockClient.capturedParent != "buildTargets/volteer" {
			t.Errorf("Parent = %q, want %q", mockClient.capturedParent, "buildTargets/volteer")
		}
		if mockClient.capturedPageSize != 50 {
			t.Errorf("PageSize = %d, want 50", mockClient.capturedPageSize)
		}

		// Verify transformed results
		expected := []string{
			"buildTargets/volteer/models/model-a",
			"buildTargets/volteer/models/model-b",
		}
		if !reflect.DeepEqual(results, expected) {
			t.Errorf("results = %v, want %v", results, expected)
		}
	})

	t.Run("returns empty slice without error when no models found", func(t *testing.T) {
		mockClient := &mockMoblabClient{
			models: []*moblabapipb.Model{},
		}

		svc := service.NewWithClient(mockClient)
		results, err := svc.ListModels(ctx, "unknown-board", 10)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(results) != 0 {
			t.Errorf("expected empty results, got len = %d", len(results))
		}
	})

	t.Run("propagates error when underlying API returns an error", func(t *testing.T) {
		expectedErr := errors.New("rpc network timeout")
		mockClient := &mockMoblabClient{
			err: expectedErr,
		}

		svc := service.NewWithClient(mockClient)
		results, err := svc.ListModels(ctx, "volteer", 10)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("error = %v, want %v", err, expectedErr)
		}
		if results != nil {
			t.Errorf("expected results to be nil on error, got %v", results)
		}
	})
}

func TestBuildServiceImpl_CheckPermission(t *testing.T) {
	ctx := context.Background()

	t.Run("returns true when model exists in list", func(t *testing.T) {
		fakeModels := []*moblabapipb.Model{
			{Name: "buildTargets/volteer/models/model-a"},
			{Name: "buildTargets/volteer/models/model-b"},
		}
		mockClient := &mockMoblabClient{
			models: fakeModels,
		}

		svc := service.NewWithClient(mockClient)
		hasPermission, err := svc.CheckPermission(ctx, "volteer", "model-b", 100)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !hasPermission {
			t.Errorf("expected permission to be true, got false")
		}
		if mockClient.capturedParent != "buildTargets/volteer" {
			t.Errorf("capturedParent = %q, want %q", mockClient.capturedParent, "buildTargets/volteer")
		}
		if mockClient.capturedPageSize != 100 {
			t.Errorf("capturedPageSize = %d, want 100", mockClient.capturedPageSize)
		}
	})

	t.Run("returns false when model is not in list", func(t *testing.T) {
		fakeModels := []*moblabapipb.Model{
			{Name: "buildTargets/volteer/models/model-a"},
		}
		mockClient := &mockMoblabClient{
			models: fakeModels,
		}

		svc := service.NewWithClient(mockClient)
		hasPermission, err := svc.CheckPermission(ctx, "volteer", "nonexistent-model", 100)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hasPermission {
			t.Errorf("expected permission to be false, got true")
		}
	})

	t.Run("returns false when model list is empty", func(t *testing.T) {
		mockClient := &mockMoblabClient{
			models: []*moblabapipb.Model{},
		}

		svc := service.NewWithClient(mockClient)
		hasPermission, err := svc.CheckPermission(ctx, "volteer", "model-a", 100)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hasPermission {
			t.Errorf("expected permission to be false, got true")
		}
	})

	t.Run("returns error when list models fails", func(t *testing.T) {
		expectedErr := errors.New("network failure")
		mockClient := &mockMoblabClient{
			err: expectedErr,
		}

		svc := service.NewWithClient(mockClient)
		hasPermission, err := svc.CheckPermission(ctx, "volteer", "model-a", 100)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !errors.Is(err, expectedErr) {
			t.Errorf("error = %v, want %v", err, expectedErr)
		}
		if hasPermission {
			t.Errorf("expected permission to be false on error, got true")
		}
	})
}
