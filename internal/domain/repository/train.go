package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type TrainRepository interface {
	Get(ctx context.Context, uid string) (*model.TrainStatus, error)
	Save(ctx context.Context, st *model.TrainStatus) error
	WithSubs(ctx context.Context) ([]string, error)
}
