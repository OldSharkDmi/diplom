package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type trainRunRepo struct {
	db *pgx.Conn
}

func NewTrainRunRepo(db *pgx.Conn) repository.TrainRunRepository {
	return &trainRunRepo{db: db}
}

func (r *trainRunRepo) GetByTrainAndDate(ctx context.Context, trainID int64, date time.Time) (*model.TrainRun, error) {
	const q = `SELECT id, train_id, run_date, created_at FROM train_runs WHERE train_id=$1 AND run_date=$2`
	row := r.db.QueryRow(ctx, q, trainID, date)
	var tr model.TrainRun
	if err := row.Scan(&tr.ID, &tr.TrainID, &tr.RunDate, &tr.CreatedAt); err != nil {
		return nil, err
	}
	return &tr, nil
}

func (r *trainRunRepo) CreateIfNotExists(ctx context.Context, trainID int64, date time.Time) (*model.TrainRun, error) {
	const q = `INSERT INTO train_runs (train_id, run_date) VALUES ($1, $2)
               ON CONFLICT (train_id, run_date) DO UPDATE SET train_id=EXCLUDED.train_id
               RETURNING id, train_id, run_date, created_at`
	row := r.db.QueryRow(ctx, q, trainID, date)
	var tr model.TrainRun
	if err := row.Scan(&tr.ID, &tr.TrainID, &tr.RunDate, &tr.CreatedAt); err != nil {
		return nil, err
	}
	return &tr, nil
}
