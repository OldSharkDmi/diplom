package model

import "time"

type Occupancy struct {
	DelayMin  int       `json:"delay_min"`
	Level     string    `json:"level"` // low|medium|high|unknown
	UpdatedAt time.Time `json:"updated_at"`
}
