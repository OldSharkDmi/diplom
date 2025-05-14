package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type Listener interface {
	Wait(ctx context.Context) (string, error)
	Close()
}

type TrainRepository interface {
	Get(ctx context.Context, uid string) (*model.TrainStatus, error)
	Save(ctx context.Context, st *model.TrainStatus) error
	WithSubs(ctx context.Context) ([]string, error)
	Occupancy(ctx context.Context, uid string) (*model.Occupancy, error)
	Listen(ctx context.Context) (Listener, error)
}

//type pgxListener struct{ *pgxpool.Conn }
