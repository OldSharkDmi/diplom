package postgres

import (
	"context"

	"train-backend/internal/domain/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type stationRepo struct {
	pool *pgxpool.Pool
}

func NewStationRepo(p *pgxpool.Pool) *stationRepo {
	return &stationRepo{p}
}

// Search ищет станции по вхождению в title, возвращает только transport_type='train'.
func (r *stationRepo) Search(ctx context.Context, q string, limit int) ([]model.Station, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT code, title, station_type, transport_type, latitude, longitude, settlement_code
          FROM stations
         WHERE title ILIKE '%' || $1 || '%'
           AND transport_type = 'train'
         ORDER BY title
         LIMIT $2
    `, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []model.Station
	for rows.Next() {
		var s model.Station
		if err := rows.Scan(
			&s.Code,
			&s.Title,
			&s.Type,
			&s.Transport,
			&s.Latitude,
			&s.Longitude,
			&s.SettlementCode,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
