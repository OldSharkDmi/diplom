package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

/* --------------------------------------------------------------------- */
/*  constructor                                                          */
/* --------------------------------------------------------------------- */

type trainRepo struct {
	pool *pgxpool.Pool
}

func NewTrainRepo(p *pgxpool.Pool) *trainRepo {
	return &trainRepo{p}
}

/* --------------------------------------------------------------------- */
/*  LISTEN / NOTIFY                                                      */
/* --------------------------------------------------------------------- */

// Listener получает уведомления из PostgreSQL-пуб/саб канала train_updates
type pgxListener struct {
	conn *pgxpool.Conn
}

// Wait ждёт следующего NOTIFY и возвращает его payload
func (l *pgxListener) Wait(ctx context.Context) (string, error) {
	// Conn() возвращает *pgx.Conn
	px := l.conn.Conn()
	// Теперь вызываем WaitForNotification у *pgx.Conn
	notif, err := px.WaitForNotification(ctx)
	if err != nil {
		return "", err
	}
	return notif.Payload, nil
}

// Close отпускает соединение обратно в пул
func (l *pgxListener) Close() {
	l.conn.Release()
}

// Listen начинает слушать канал train_updates и возвращает Listener
func (r *trainRepo) Listen(ctx context.Context) (repository.Listener, error) {
	acq, err := r.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	// регистрируемся на уведомления
	if _, err = acq.Exec(ctx, `LISTEN train_updates`); err != nil {
		acq.Release()
		return nil, err
	}
	return &pgxListener{conn: acq}, nil
}

/* ----------------------------- CRUD-методы ------------------------------ */

func (r *trainRepo) Get(ctx context.Context, uid string) (*model.TrainStatus, error) {
	var raw []byte
	err := r.pool.QueryRow(ctx,
		`SELECT last_status FROM trains WHERE uid=$1`, uid,
	).Scan(&raw)
	if err != nil {
		return nil, err
	}
	var st model.TrainStatus
	return &st, json.Unmarshal(raw, &st)
}

func (r *trainRepo) Save(ctx context.Context, st *model.TrainStatus) error {
	buf, _ := json.Marshal(st)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO trains(uid,last_status,last_update)
		 VALUES ($1,$2,now())
		 ON CONFLICT (uid) DO UPDATE
		   SET last_status = EXCLUDED.last_status,
		       last_update = now()`,
		st.UID, buf,
	)
	return err
}

func (r *trainRepo) WithSubs(ctx context.Context) ([]string, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT DISTINCT train_uid FROM subscriptions`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var uids []string
	for rows.Next() {
		var u string
		if err := rows.Scan(&u); err != nil {
			return nil, err
		}
		uids = append(uids, u)
	}
	return uids, rows.Err()
}

func (r *trainRepo) Occupancy(ctx context.Context, uid string) (*model.Occupancy, error) {
	const q = `SELECT delay_min, occupancy, updated_at FROM v_train_occupancy WHERE uid=$1`
	var o model.Occupancy
	if err := r.pool.QueryRow(ctx, q, uid).Scan(&o.DelayMin, &o.Level, &o.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) { // ← 204
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}
