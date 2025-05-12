package repository

import (
	"context"
	"time"
	"train-backend/internal/domain/model"
)

type TrainRunRepository interface {
	GetByTrainAndDate(ctx context.Context, trainID int64, date time.Time) (*model.TrainRun, error)
	CreateIfNotExists(ctx context.Context, trainID int64, date time.Time) (*model.TrainRun, error)
}
