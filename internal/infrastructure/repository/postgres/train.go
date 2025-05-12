package postgres

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/domain/model"
)

type trainRepo struct{ pool *pgxpool.Pool }

func NewTrainRepo(p *pgxpool.Pool) *trainRepo { return &trainRepo{p} }

/* Get stored status */
func (r *trainRepo) Get(ctx context.Context, uid string) (*model.TrainStatus, error) {
	var raw []byte
	if err := r.pool.QueryRow(ctx,
		`SELECT last_status FROM trains WHERE uid=$1`, uid).Scan(&raw); err != nil {
		return nil, err
	}
	var st model.TrainStatus
	return &st, json.Unmarshal(raw, &st)
}

/* Upsert new status */
func (r *trainRepo) Save(ctx context.Context, st *model.TrainStatus) error {
	buf, _ := json.Marshal(st)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO trains(uid,last_status,last_update)
         VALUES ($1,$2,now())
         ON CONFLICT (uid) DO
         UPDATE SET last_status = EXCLUDED.last_status,
                    last_update = now()`,
		st.UID, buf)
	return err
}

/* Return every uid that has active subscriptions */
func (r *trainRepo) WithSubs(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT train_uid FROM subscriptions`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var uids []string
	for rows.Next() {
		var u string
		_ = rows.Scan(&u)
		uids = append(uids, u)
	}
	return uids, rows.Err()
}

func (r *trainRepo) Occupancy(ctx context.Context, uid string) (*model.Occupancy, error) {
	const q = `SELECT delay_min, occupancy, updated_at FROM v_train_occupancy WHERE uid=$1`
	var o model.Occupancy
	return &o, r.pool.QueryRow(ctx, q, uid).Scan(&o.DelayMin, &o.Level, &o.UpdatedAt)
}
