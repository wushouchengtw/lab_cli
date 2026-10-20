package service

import (
	"context"
	"fmt"
	"net/http"

	"go.chromium.org/luci/auth"
	bbgrpcpb "go.chromium.org/luci/buildbucket/proto/grpcpb"
	"go.chromium.org/luci/grpc/prpc"
	"go.chromium.org/luci/hardcoded/chromeinfra"
)

const BuildbucketHost = "cr-buildbucket.appspot.com"

// NewBuildBucketClient creates an authenticated BuildsClient using a service account JSON file.
// If serviceAccountPath is empty, it falls back to ambient credentials (e.g., local login).
func NewBuildBucketClient(ctx context.Context, serviceAccountPath string) (bbgrpcpb.BuildsClient, error) {
	opts := chromeinfra.SetDefaultAuthOptions(auth.Options{
		ServiceAccountJSONPath: serviceAccountPath,
	})
	loginMode := auth.SilentLogin
	if serviceAccountPath == "" {
		loginMode = auth.OptionalLogin
	}
	httpClient, err := auth.NewAuthenticator(ctx, loginMode, opts).Client()
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticated http client: %w", err)
	}
	return NewBuildBucketClientWithHTTPClient(httpClient), nil
}

// NewBuildBucketClientWithHTTPClient creates a BuildsClient from a custom http.Client.
func NewBuildBucketClientWithHTTPClient(httpClient *http.Client) bbgrpcpb.BuildsClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	prpcClient := &prpc.Client{
		Host: BuildbucketHost,
		C:    httpClient,
	}
	return bbgrpcpb.NewBuildsClient(prpcClient)
}
