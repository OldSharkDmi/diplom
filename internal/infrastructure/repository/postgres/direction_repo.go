package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type directionRepo struct {
	pool *pgxpool.Pool
}

func NewDirectionRepo(p *pgxpool.Pool) *directionRepo { return &directionRepo{p} }

func (r *directionRepo) Fetch(ctx context.Context, offset, limit int) ([]model.Direction, int, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, from_station, to_station
		   FROM directions
		   ORDER BY id
		   OFFSET $1 LIMIT $2`, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var dirs []model.Direction
	for rows.Next() {
		var d model.Direction
		if err = rows.Scan(&d.ID, &d.Name, &d.From, &d.To); err != nil {
			return nil, 0, err
		}
		dirs = append(dirs, d)
	}
	var total int
	r.pool.QueryRow(ctx, `SELECT count(*) FROM directions`).Scan(&total)
	return dirs, total, rows.Err()
}
