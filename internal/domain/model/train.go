package model

import (
	"encoding/json"
	"time"
)

type TrainStatus struct {
	ID         int64           `json:"id"`
	TrainRunID int64           `json:"train_run_id"`
	Status     string          `json:"status"`
	ReceivedAt time.Time       `json:"received_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Raw        json.RawMessage `json:"raw,omitempty"`
}
