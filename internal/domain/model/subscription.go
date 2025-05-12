package model

type Subscription struct {
	ID          int64  `json:"id"`
	DeviceToken string `json:"device_token"`
	TrainUID    string `json:"train_uid"`
}
