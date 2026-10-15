package service

import (
	"context"
	"fmt"
	"log/slog"
	"path"
)

func ParseBuildTargetsBoardPath(board string) string {
	return fmt.Sprintf("buildTargets/%s", board)
}

func GetModelName(model string) string {
	return path.Base(model)
}

func (b *BuildServiceImpl) ListModels(ctx context.Context, board string, pageSize int) ([]string, error) {
	slog.Debug("Listing models for board", "board", board)

	parent := ParseBuildTargetsBoardPath(board)

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

func (s *BuildServiceImpl) CheckPermission(ctx context.Context, board string, model string, limit int) (bool, error) {
	models, err := s.ListModels(ctx, board, limit)
	if err != nil {
		return false, fmt.Errorf("failed to list models for board %q: %w", board, err)
	}

	slog.Debug("fetched models successfully", "board", board, "count", len(models), "models", models)

	for _, m := range models {
		if model == GetModelName(m) {
			return true, nil
		}
	}
	return false, nil
}
