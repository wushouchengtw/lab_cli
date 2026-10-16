package buildbucket

import (
	"net/http"

	bbgrpcpb "go.chromium.org/luci/buildbucket/proto/grpcpb"
	"go.chromium.org/luci/grpc/prpc"
)

const BuildbucketHost = "cr-buildbucket.appspot.com"

func NewClient(httpClient *http.Client) bbgrpcpb.BuildsClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	prpcClient := &prpc.Client{
		Host: BuildbucketHost,
		C:    httpClient,
	}
	return bbgrpcpb.NewBuildsClient(prpcClient)
}
