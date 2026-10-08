package service

import (
	"context"
	"fmt"
	"lab_cli/site"

	moblabapi "go.chromium.org/infra/cros/satlab/common/google.golang.org/google/chromeos/moblab"
	"google.golang.org/api/option"
)

type BuildServiceImpl struct {
	// client the `BuildClient`
	client *moblabapi.BuildClient
}

// New sets up the `BuildClient` and returns a BuildConnector.
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
	return &BuildServiceImpl{
		client: client,
	}, nil
}

type IBuildService interface {
	// ListModels returns all models by given board.
	ListModels(ctx context.Context, board string, pageSize int) ([]string, error)
}
