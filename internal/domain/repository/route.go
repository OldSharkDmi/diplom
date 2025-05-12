package repository

import (
	"context"
	"time"

	"train-backend/internal/domain/model"
)

type RouteRepository interface {
	SaveCache(ctx context.Context, from, to string, date time.Time, r *model.Route) error
	GetCache(ctx context.Context, from, to string, date time.Time) (*model.Route, error)
	ByID(ctx context.Context, id int64) (*model.Route, error)
}
