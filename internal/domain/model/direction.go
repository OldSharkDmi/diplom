package model

type Direction struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
}
