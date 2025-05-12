package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type TrainRepository interface {
	Get(ctx context.Context, uid string) (*model.TrainStatus, error)
	Save(ctx context.Context, st *model.TrainStatus) error
	WithSubs(ctx context.Context) ([]string, error)
	Occupancy(ctx context.Context, uid string) (*model.Occupancy, error) // ← было OccupancyPrediction
	Listen(ctx context.Context) (*pgxListener, error)                    // для WS
}

func (r *trainRepo) Listen(ctx context.Context) (*pgx.Listener, error) {
	l, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	_, _ = l.Exec(ctx, `LISTEN train_updates`)
	return &pgxListener{Conn: l}, nil
}

type pgxListener struct{ *pgxpool.Conn }

func (l *pgxListener) WaitForNotification(ctx context.Context) (string, error) {
	n, err := l.Conn.Conn().WaitForNotification(ctx)
	if err != nil {
		return "", err
	}
	return n.Payload, nil
}
