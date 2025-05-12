package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type subRepo struct{ pool *pgxpool.Pool }

func NewSubRepo(p *pgxpool.Pool) *subRepo { return &subRepo{p} }

func (r *subRepo) Create(ctx context.Context, s *model.Subscription) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO subscriptions (device_token,train_uid)
         VALUES ($1,$2) RETURNING id`,
		s.DeviceToken, s.TrainUID).Scan(&s.ID)
}

func (r *subRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM subscriptions WHERE id=$1`, id)
	return err
}

func (r *subRepo) ByTrain(ctx context.Context, uid string) ([]model.Subscription, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id,device_token,train_uid FROM subscriptions WHERE train_uid=$1`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Subscription
	for rows.Next() {
		var s model.Subscription
		_ = rows.Scan(&s.ID, &s.DeviceToken, &s.TrainUID)
		out = append(out, s)
	}
	return out, rows.Err()
}
