package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type SubRepository interface {
	Create(ctx context.Context, s *model.Subscription) error
	Delete(ctx context.Context, id int64) error
	ByTrain(ctx context.Context, uid string) ([]model.Subscription, error)
}
