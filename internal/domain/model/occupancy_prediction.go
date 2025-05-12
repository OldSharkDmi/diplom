package model

import "time"

type OccupancyPrediction struct {
	ID          int64     `json:"id"`
	TrainRunID  int64     `json:"train_run_id"`
	CarNumber   int16     `json:"car_number"`
	Level       string    `json:"level"` // low/medium/high
	PredictedAt time.Time `json:"predicted_at"`
}
