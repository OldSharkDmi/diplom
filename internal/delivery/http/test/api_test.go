package http_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080/api/v1"

func TestHealth(t *testing.T) {
	res, err := http.Get(baseURL + "/health")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 200, res.StatusCode)

	var h struct {
		Status string `json:"status,omitempty"`
		DB     string `json:"db,omitempty"`
		Cache  string `json:"cache,omitempty"`
	}
	err = json.NewDecoder(res.Body).Decode(&h)
	require.NoError(t, err)
	require.Equal(t, "ok", h.Status)
}

func TestStationsSearch(t *testing.T) {
	res, err := http.Get(baseURL + "/stations?search=одинц&limit=5")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 200, res.StatusCode)

	var out struct {
		Data []struct {
			Code          string  `json:"code"`
			Title         string  `json:"title"`
			TransportType string  `json:"transport_type"`
			Latitude      float64 `json:"latitude"`
			Longitude     float64 `json:"longitude"`
		} `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&out)
	require.NoError(t, err)
	require.NotEmpty(t, out.Data)
}

func TestStationSchedule(t *testing.T) {
	code := "s9600721"
	date := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	res, err := http.Get(baseURL + "/station/" + code + "?date=" + date + "&limit=5")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 200, res.StatusCode)

	var out struct {
		Date       string                             `json:"date"`
		Pagination struct{ Total, Limit, Offset int } `json:"pagination"`
		Schedule   []interface{}                      `json:"schedule"`
	}
	err = json.NewDecoder(res.Body).Decode(&out)
	require.NoError(t, err)
	require.Equal(t, date, out.Date)
	require.GreaterOrEqual(t, out.Pagination.Limit, 1)
}

func TestDirectSearch(t *testing.T) {
	from, to := "s9600721", "s2000006"
	date := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	url := baseURL + "/search?from=" + from + "&to=" + to + "&date=" + date + "&limit=3"
	res, err := http.Get(url)
	require.NoError(t, err)
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadGateway {
		t.Log("Пропускаем прямой поиск — Яндекс вернул 502.")
		return
	}
	require.Equal(t, 200, res.StatusCode)

	var out struct {
		Pagination struct{ Total, Limit, Offset int } `json:"pagination"`
		Segments   []interface{}                      `json:"segments"`
	}
	err = json.NewDecoder(res.Body).Decode(&out)
	require.NoError(t, err)
}

func TestRouteWithTransfer(t *testing.T) {
	from, to := "s9600721", "s2000006"
	date := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	url := baseURL + "/routes?from=" + from + "&to=" + to + "&date=" + date
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 200, res.StatusCode)

	var out struct {
		Duration float64       `json:"duration"`
		Segments []interface{} `json:"segments"`
	}
	err = json.NewDecoder(res.Body).Decode(&out)
	require.NoError(t, err)
}

func TestGetTrainStatus_NotFound(t *testing.T) {
	// случайный UID, которого точно нет
	uid := "no_such_train_12345"
	res, err := http.Get(baseURL + "/trains/" + uid)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 404, res.StatusCode)
}

func TestGetTrainOccupancy_NotFound(t *testing.T) {
	uid := "no_such_train_12345"
	res, err := http.Get(baseURL + "/trains/" + uid + "/occupancy")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 404, res.StatusCode)
}

func TestSubscriptions_CRUD(t *testing.T) {
	// 1. Create
	payload := map[string]string{"device_token": "tok123", "train_uid": "abc123"}
	body, _ := json.Marshal(payload)
	res, err := http.Post(baseURL+"/subscriptions", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 201, res.StatusCode)

	var sub struct {
		ID          int64  `json:"id"`
		DeviceToken string `json:"device_token"`
		TrainUID    string `json:"train_uid"`
	}
	err = json.NewDecoder(res.Body).Decode(&sub)
	require.NoError(t, err)
	require.Equal(t, "tok123", sub.DeviceToken)
	require.Equal(t, "abc123", sub.TrainUID)

	// 2. Delete
	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/subscriptions/"+strconv.FormatInt(sub.ID, 10), nil)
	res2, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res2.Body.Close()
	require.Equal(t, 204, res2.StatusCode)
}

func TestEvents(t *testing.T) {
	payload := map[string]interface{}{
		"device_id":  "dev42",
		"event_type": "open_screen",
		"payload":    map[string]string{"screen": "search"},
	}
	body, _ := json.Marshal(payload)
	res, err := http.Post(baseURL+"/events", "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, 201, res.StatusCode)

	var ev struct {
		ID        int64       `json:"id"`
		DeviceID  string      `json:"device_id"`
		EventType string      `json:"event_type"`
		Payload   interface{} `json:"payload"`
	}
	err = json.NewDecoder(res.Body).Decode(&ev)
	require.NoError(t, err)
	require.Equal(t, "dev42", ev.DeviceID)
	require.Equal(t, "open_screen", ev.EventType)
}
