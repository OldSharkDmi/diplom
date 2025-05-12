package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type occupancyRepo struct {
	db *pgx.Conn
}

func NewOccupancyRepo(db *pgx.Conn) repository.OccupancyPredictionRepository {
	return &occupancyRepo{db: db}
}

func (r *occupancyRepo) GetLatestByRunID(ctx context.Context, runID int64) ([]*model.OccupancyPrediction, error) {
	const q = `SELECT id, train_run_id, car_number, level, predicted_at
               FROM occupancy_predictions
               WHERE train_run_id=$1
               ORDER BY predicted_at DESC`
	rows, err := r.db.Query(ctx, q, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var res []*model.OccupancyPrediction
	for rows.Next() {
		var op model.OccupancyPrediction
		if err := rows.Scan(&op.ID, &op.TrainRunID, &op.CarNumber, &op.Level, &op.PredictedAt); err != nil {
			return nil, err
		}
		res = append(res, &op)
	}
	return res, nil
}

func (r *occupancyRepo) SaveBatch(ctx context.Context, preds []*model.OccupancyPrediction) error {
	batch := &pgx.Batch{}
	const q = `INSERT INTO occupancy_predictions (train_run_id, car_number, level, predicted_at)
               VALUES ($1, $2, $3, $4)
               ON CONFLICT (train_run_id, car_number) DO UPDATE SET level=EXCLUDED.level, predicted_at=EXCLUDED.predicted_at`
	for _, p := range preds {
		batch.Queue(q, p.TrainRunID, p.CarNumber, p.Level, p.PredictedAt)
	}
	br := r.db.SendBatch(ctx, batch)
	return br.Close()
}
