package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"lab_cli/service"

	"google.golang.org/api/iterator"
	moblabapipb "google.golang.org/genproto/googleapis/chromeos/moblab/v1beta1"
)

// mockModelIterator mocks the Google API iterator behavior.
type mockModelIterator struct {
	items []*moblabapipb.Model
	err   error
	index int
}

func (m *mockModelIterator) Next() (*moblabapipb.Model, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.index >= len(m.items) {
		return nil, iterator.Done
	}
	item := m.items[m.index]
	m.index++
	return item, nil
}

// mockMoblabClient mocks the client and captures the request for assertion.
type mockMoblabClient struct {
	iter        service.IModelIterator
	capturedReq *moblabapipb.ListModelsRequest
}

func (m *mockMoblabClient) ListModels(ctx context.Context, req *moblabapipb.ListModelsRequest) service.IModelIterator {
	m.capturedReq = req
	return m.iter
}

func TestParseBuildTargetsPath(t *testing.T) {
	got := service.ParseBuildTargetsPath("volteer")
	want := "buildTargets/volteer"
	if got != want {
		t.Errorf("ParseBuildTargetsPath() = %q, want %q", got, want)
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
			iter: &mockModelIterator{items: fakeModels},
		}

		svc := service.NewWithClient(mockClient)
		results, err := svc.ListModels(ctx, "volteer", 50)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Verify request parameters are properly constructed
		if mockClient.capturedReq == nil {
			t.Fatal("expected ListModelsRequest to be captured, got nil")
		}
		if mockClient.capturedReq.Parent != "buildTargets/volteer" {
			t.Errorf("Parent = %q, want %q", mockClient.capturedReq.Parent, "buildTargets/volteer")
		}
		if mockClient.capturedReq.PageSize != 50 {
			t.Errorf("PageSize = %d, want 50", mockClient.capturedReq.PageSize)
		}

		// Verify parsed results
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
			iter: &mockModelIterator{items: []*moblabapipb.Model{}},
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
			iter: &mockModelIterator{err: expectedErr},
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
