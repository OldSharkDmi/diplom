package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type stationRepo struct{ pool *pgxpool.Pool }

func NewStationRepo(p *pgxpool.Pool) *stationRepo { return &stationRepo{p} }

func (r *stationRepo) Search(ctx context.Context, q string, limit int) ([]model.Station, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT code,title,station_type,transport_type,latitude,longitude,settlement_code
		   FROM stations
		  WHERE transport_type='suburban'
		    AND title ILIKE '%'||$1||'%'
		  ORDER BY similarity(title,$1) DESC
		  LIMIT $2`, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Station
	for rows.Next() {
		var s model.Station
		if err = rows.Scan(&s.Code, &s.Title, &s.Type, &s.Transport,
			&s.Latitude, &s.Longitude, &s.SettlementCode); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
