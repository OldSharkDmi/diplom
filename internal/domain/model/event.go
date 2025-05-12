package model

type Event struct {
	ID       int64       `json:"id"`
	Ts       string      `json:"ts"`
	DeviceID string      `json:"device_id,omitempty"`
	Type     string      `json:"event_type"`
	Payload  interface{} `json:"payload"`
}
