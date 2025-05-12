package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type TrainStatusRepository interface {
	GetByRunID(ctx context.Context, runID int64) (*model.TrainStatus, error)
	Upsert(ctx context.Context, status *model.TrainStatus) error
}
