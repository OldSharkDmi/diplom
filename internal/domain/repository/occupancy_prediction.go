package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type OccupancyPredictionRepository interface {
	GetLatestByRunID(ctx context.Context, runID int64) ([]*model.OccupancyPrediction, error)
	SaveBatch(ctx context.Context, preds []*model.OccupancyPrediction) error
}
