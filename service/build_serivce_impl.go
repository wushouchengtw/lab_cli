package service

import (
	"context"
	"fmt"
	"lab_cli/models"
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

func (s *BuildServiceImpl) CheckPermission(ctx context.Context, board string, modelsOnCheck []string, limit int) ([]models.CheckPermissionResponse, error) {
	modelsWithPermission, err := s.ListModels(ctx, board, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list models for board %q: %w", board, err)
	}

	slog.Debug("fetched models successfully", "board", board, "count", len(modelsWithPermission), "models", modelsWithPermission)

	result := []models.CheckPermissionResponse{}
	modelsPermissionMap := make(map[string]bool)
	for _, model := range modelsWithPermission {
		modelsPermissionMap[GetModelName(model)] = true
	}

	// Check permission for each model in modelsOnCheck.
	for _, model := range modelsOnCheck {
		result = append(result, models.CheckPermissionResponse{
			Board:         board,
			Model:         model,
			HasPermission: modelsPermissionMap[model],
		})
	}
	return result, nil
}
