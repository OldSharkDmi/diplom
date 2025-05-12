package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type routeRepo struct{ pool *pgxpool.Pool }

func NewRouteRepo(p *pgxpool.Pool) *routeRepo { return &routeRepo{p} }

func (r *routeRepo) GetCache(ctx context.Context, from, to string, date time.Time) (*model.Route, error) {
	var raw []byte
	err := r.pool.QueryRow(ctx,
		`SELECT json FROM routes
          WHERE from_code=$1 AND to_code=$2 AND date=$3
          AND cached_at > now() - interval '6 hours'`,
		from, to, date).Scan(&raw)
	if err != nil {
		return nil, err
	}
	var rt model.Route
	return &rt, json.Unmarshal(raw, &rt)
}

func (r *routeRepo) SaveCache(ctx context.Context, from, to string, date time.Time, rt *model.Route) error {
	buf, _ := json.Marshal(rt)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO routes (from_code,to_code,date,json)
         VALUES ($1,$2,$3,$4)
         ON CONFLICT (from_code,to_code,date) DO
         UPDATE SET json=EXCLUDED.json, cached_at=now()`,
		from, to, date, buf)
	return err
}
