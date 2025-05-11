package model

type Station struct {
	Code           string  `json:"code"`
	Title          string  `json:"title"`
	Type           string  `json:"station_type"`
	Transport      string  `json:"transport_type"`
	Latitude       float64 `json:"latitude,omitempty"`
	Longitude      float64 `json:"longitude,omitempty"`
	SettlementCode string  `json:"settlement_code,omitempty"`
}
