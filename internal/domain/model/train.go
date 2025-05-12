package model

type TrainStatus struct {
	UID        string `json:"uid"`
	Departure  string `json:"departure,omitempty"`
	Arrival    string `json:"arrival,omitempty"`
	DelayMin   int    `json:"delay_min,omitempty"`
	Occupancy  string `json:"occupancy,omitempty"` // low / medium / high
	UpdatedISO string `json:"updated_at"`
}
