package service

import (
	"context"
	"fmt"

	"go.chromium.org/infra/cros/satlab/common/utils/collection"
	moblabapipb "google.golang.org/genproto/googleapis/chromeos/moblab/v1beta1"
)

func ParseBuildTargetsPath(board string) string {
	return fmt.Sprintf("buildTargets/%s", board)
}

func (b *BuildServiceImpl) ListModels(ctx context.Context, board string, pageSize int) ([]string, error) {
	fmt.Println("Listing models for board: ", board)

	parent := ParseBuildTargetsPath(board)

	req := &moblabapipb.ListModelsRequest{
		Parent:   parent,
		PageSize: int32(pageSize),
	}

	iter := b.client.ListModels(ctx, req)

	res, err := collection.Collect(
		iter.Next,
		func(model *moblabapipb.Model) (string, error) {
			return model.GetName(), nil
		},
	)

	if err != nil {
		return nil, err
	}

	return res, nil
}
