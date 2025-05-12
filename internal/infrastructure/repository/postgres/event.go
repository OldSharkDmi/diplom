package postgres

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type eventRepo struct{ pool *pgxpool.Pool }

func NewEventRepo(p *pgxpool.Pool) *eventRepo { return &eventRepo{p} }

func (r *eventRepo) Store(ctx context.Context, e *model.Event) error {
	payload, _ := json.Marshal(e.Payload)
	return r.pool.QueryRow(ctx,
		`INSERT INTO events(device_id,event_type,payload)
         VALUES ($1,$2,$3) RETURNING id,ts`,
		e.DeviceID, e.Type, payload).
		Scan(&e.ID, &e.Ts)
}
