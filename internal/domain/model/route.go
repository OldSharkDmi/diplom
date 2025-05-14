package model

type RouteSegment struct {
	FromCode string `json:"from"`
	ToCode   string `json:"to"`
	Dep      string `json:"dep"`
	Arr      string `json:"arr"`
	TrainUID string `json:"train_uid"`
}

type Route struct {
	Segments []RouteSegment `json:"segments"`
	Duration float64        `json:"duration"`
}
