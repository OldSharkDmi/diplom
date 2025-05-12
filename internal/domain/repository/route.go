package repository

import (
	"context"
	"time"

	"train-backend/internal/domain/model"
)

type RouteRepository interface {
	GetCache(ctx context.Context, from, to string, d time.Time) (*model.Route, error)
	SaveCache(ctx context.Context, from, to string, d time.Time, r *model.Route) error
}
