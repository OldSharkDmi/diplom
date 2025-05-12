package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type trainStatusRepo struct {
	db *pgx.Conn
}

func NewTrainStatusRepo(db *pgx.Conn) repository.TrainStatusRepository {
	return &trainStatusRepo{db: db}
}

func (r *trainStatusRepo) GetByRunID(ctx context.Context, runID int64) (*model.TrainStatus, error) {
	const q = `SELECT id, train_run_id, status, received_at, updated_at, raw FROM train_statuses WHERE train_run_id=$1`
	row := r.db.QueryRow(ctx, q, runID)
	var ts model.TrainStatus
	if err := row.Scan(&ts.ID, &ts.TrainRunID, &ts.Status, &ts.ReceivedAt, &ts.UpdatedAt, &ts.Raw); err != nil {
		return nil, err
	}
	return &ts, nil
}

func (r *trainStatusRepo) Upsert(ctx context.Context, status *model.TrainStatus) error {
	const q = `INSERT INTO train_statuses (train_run_id, status, received_at, updated_at, raw)
               VALUES ($1, $2, $3, $4, $5)
               ON CONFLICT (train_run_id) DO UPDATE SET status=EXCLUDED.status, received_at=EXCLUDED.received_at, updated_at=EXCLUDED.updated_at, raw=EXCLUDED.raw`
	_, err := r.db.Exec(ctx, q, status.TrainRunID, status.Status, status.ReceivedAt, status.UpdatedAt, status.Raw)
	return err
}
