package service

import (
	"context"
	"fmt"
)

func ParseBuildTargetsPath(board string) string {
	return fmt.Sprintf("buildTargets/%s", board)
}

func (b *BuildServiceImpl) ListModels(ctx context.Context, board string, pageSize int) ([]string, error) {
	fmt.Println("Listing models for board: ", board)

	parent := ParseBuildTargetsPath(board)

	models, err := b.client.ListModels(ctx, parent, int32(pageSize))
	if err != nil {
		return nil, err
	}

	res := make([]string, len(models))
	for i, m := range models {
		res[i] = m.GetName()
	}

	return res, nil
}
