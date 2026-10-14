package service

import (
	"context"
	"fmt"
	"lab_cli/site"

	moblabapi "go.chromium.org/infra/cros/satlab/common/google.golang.org/google/chromeos/moblab"
	"google.golang.org/api/option"
	moblabapipb "google.golang.org/genproto/googleapis/chromeos/moblab/v1beta1"
)

type IModelIterator interface {
	Next() (*moblabapipb.Model, error)
}

type IMoblabClient interface {
	ListModels(ctx context.Context, req *moblabapipb.ListModelsRequest) IModelIterator
}

type moblabClientAdapter struct {
	client *moblabapi.BuildClient
}

func (a *moblabClientAdapter) ListModels(ctx context.Context, req *moblabapipb.ListModelsRequest) IModelIterator {
	return a.client.ListModels(ctx, req)
}

type BuildServiceImpl struct {
	// client the `IMoblabClient` interface
	client IMoblabClient
}

// New sets up the `BuildClient` and returns an IBuildService.
// The service account is set in the global environment.
func New(ctx context.Context, serviceAccountPath ...string) (IBuildService, error) {
	// create moblab client using service account json file
	path := site.GetServiceAccountPath(serviceAccountPath...)
	if path == "" {
		return nil, fmt.Errorf("service account path is not set")
	}
	client, err := moblabapi.NewBuildClient(ctx, option.WithCredentialsFile(path))
	if err != nil {
		return nil, err
	}
	return NewWithClient(&moblabClientAdapter{client: client}), nil
}

// NewWithClient creates an IBuildService with an injected IMoblabClient.
func NewWithClient(client IMoblabClient) IBuildService {
	return &BuildServiceImpl{
		client: client,
	}
}

type IBuildService interface {
	// ListModels returns all models by given board.
	ListModels(ctx context.Context, board string, pageSize int) ([]string, error)
}
