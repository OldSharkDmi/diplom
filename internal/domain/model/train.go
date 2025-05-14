package model

import (
	"encoding/json"
	"time"
)

type TrainStatus struct {
	/* ключевые поля для API */
	UID string `json:"uid"` // TEST_001

	Departure  string `json:"departure"`   // 2025-05-15T08:32
	Arrival    string `json:"arrival"`     // 2025-05-15T10:54
	UpdatedISO string `json:"updated_iso"` // RFC-3339
	DelayMin   int    `json:"delay_min"`   // –2…+60
	Occupancy  string `json:"occupancy"`   // low|medium|high|unknown

	/* ← добавляем обратно */
	Status string `json:"status,omitempty"` // on_time|delayed|cancelled

	/* служебные поля */
	ID         int64           `json:"id,omitempty"`
	TrainRunID int64           `json:"train_run_id,omitempty"`
	ReceivedAt time.Time       `json:"received_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Raw        json.RawMessage `json:"raw,omitempty"`
}
