package handler

import "time"

type TrainRun struct {
	ID        int64     `json:"id"`
	TrainID   int64     `json:"train_id"`
	RunDate   time.Time `json:"run_date"`
	CreatedAt time.Time `json:"created_at"`
}
