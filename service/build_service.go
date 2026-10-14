package service

import (
	"context"
	"fmt"
	"lab_cli/site"

	moblabapi "go.chromium.org/infra/cros/satlab/common/google.golang.org/google/chromeos/moblab"
	"go.chromium.org/infra/cros/satlab/common/utils/collection"
	"google.golang.org/api/option"
	moblabapipb "google.golang.org/genproto/googleapis/chromeos/moblab/v1beta1"
)

// IMoblabClient defines the minimal interface for interacting with Moblab models.
type IMoblabClient interface {
	ListModels(ctx context.Context, parent string, pageSize int32) ([]*moblabapipb.Model, error)
}

// moblabClientAdapter encapsulates the Google SDK Iterator consumption.
type moblabClientAdapter struct {
	client *moblabapi.BuildClient
}

func (a *moblabClientAdapter) ListModels(ctx context.Context, parent string, pageSize int32) ([]*moblabapipb.Model, error) {
	req := &moblabapipb.ListModelsRequest{
		Parent:   parent,
		PageSize: pageSize,
	}
	iter := a.client.ListModels(ctx, req)
	return collection.Collect(
		iter.Next,
		func(m *moblabapipb.Model) (*moblabapipb.Model, error) {
			return m, nil
		},
	)
}

type BuildServiceImpl struct {
	client IMoblabClient
}

// New sets up the BuildClient and returns a *BuildServiceImpl.
func New(ctx context.Context, serviceAccountPath ...string) (*BuildServiceImpl, error) {
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

// NewWithClient creates a *BuildServiceImpl with an injected IMoblabClient.
func NewWithClient(client IMoblabClient) *BuildServiceImpl {
	return &BuildServiceImpl{
		client: client,
	}
}
