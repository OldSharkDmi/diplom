package yandex

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const baseURL = "https://api.rasp.yandex.net/v3.0"

// DTO для /thread/
type ThreadResponse struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	Stops     []Stop `json:"stops"`
}
type Stop struct {
	Station Station `json:"station"`
	Dep     string  `json:"departure"`
	Arr     string  `json:"arrival"`
}

func (c *Client) Thread(ctx context.Context, uid, date string) (*ThreadResponse, error) {
	q := url.Values{
		"uid":  {uid},
		"date": {date},
	}
	var out ThreadResponse
	return &out, c.do(ctx, "thread/", q, &out)
}

// Client — публичный тип
type Client struct {
	apiKey string
	http   *http.Client
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}
func (c *Client) SetTimeout(d time.Duration) { c.http.Timeout = d }

// запрос с общей обработкой ошибок
func (c *Client) do(ctx context.Context, endpoint string, q url.Values, out any) error {
	q.Set("apikey", c.apiKey)
	u := fmt.Sprintf("%s/%s?%s", baseURL, endpoint, q.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	log.Printf("[Yandex] %s %s", req.Method, req.URL.String())
	res, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		msg, _ := io.ReadAll(io.LimitReader(res.Body, 4<<10))
		return fmt.Errorf("yandex api %s: %s", res.Status, msg)
	}
	return json.NewDecoder(res.Body).Decode(out)
}
func (c *Client) StationsList(ctx context.Context, transport []string) ([]Station, error) {
	q := url.Values{"lang": {"ru_RU"}}
	if len(transport) > 0 {
		q.Set("transport_types", strings.Join(transport, ","))
	}

	var raw struct {
		Countries []struct {
			Regions []struct {
				Settlements []struct {
					Code     string    `json:"code"`
					Stations []Station `json:"stations"`
				} `json:"settlements"`
			} `json:"regions"`
		} `json:"countries"`
	}
	if err := c.do(ctx, "stations_list/", q, &raw); err != nil {
		return nil, err
	}

	var out []Station
	for _, ctry := range raw.Countries {
		for _, reg := range ctry.Regions {
			for _, set := range reg.Settlements {
				for _, st := range set.Stations {
					st.Code = st.Codes.Yandex
					if st.Code == "" { // пропускаем безкодовые платформы
						continue
					}
					st.SettlementCode = set.Code
					out = append(out, st) // ← просто добавляем
				}

			}
		}
	}
	return out, nil
}

/* ───── DTO ───── */

type (
	Station struct {
		Codes struct {
			Yandex string `json:"yandex_code"`
		} `json:"codes"`

		// теперь из JSON будет правильно браться поле "code"
		Code           string `json:"code"`
		Type           string `json:"station_type"`
		Title          string `json:"title"`
		Transport      string `json:"transport_type"`
		Latitude       Num    `json:"latitude"`
		Longitude      Num    `json:"longitude"`
		SettlementCode string `json:"-"`
	}
	Thread struct {
		UID       string `json:"uid"`
		Title     string `json:"title"`
		Number    string `json:"number"`
		Transport string `json:"transport_type"`
	}

	Segment struct {
		Thread    Thread  `json:"thread"`
		Departure string  `json:"departure"`
		Arrival   string  `json:"arrival"`
		Duration  float64 `json:"duration"`
		From      Station `json:"from"`
		To        Station `json:"to"`
	}
	SearchResponse struct {
		Pagination struct {
			Total  int `json:"total"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		} `json:"pagination"`
		Segments []Segment `json:"segments"`
	}
	ScheduleResponse struct {
		Date       string `json:"date"`
		Pagination struct {
			Total  int `json:"total"`
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		} `json:"pagination"`
		Station  Station   `json:"station"`
		Schedule []Segment `json:"schedule"`
	}
)

/* ───── endpoints ───── */

// /search/
func (c *Client) Search(ctx context.Context, from, to, date string,
	transport []string, transfers bool, offset, limit int) (*SearchResponse, error) {

	q := url.Values{
		"from":   {from},
		"to":     {to},
		"date":   {date},
		"offset": {strconv.Itoa(offset)},
		"limit":  {strconv.Itoa(limit)},
	}
	if len(transport) > 0 {
		q.Set("transport_types", strings.Join(transport, ","))
	}
	if transfers {
		q.Set("transfers", "true")
	}

	var resp SearchResponse
	return &resp, c.do(ctx, "search/", q, &resp)
}

// /schedule/
func (c *Client) ScheduleOnStation(ctx context.Context, station, date, event string,
	transport []string, offset, limit int) (*ScheduleResponse, error) {

	q := url.Values{
		"station": {station},
		"date":    {date},
		"offset":  {strconv.Itoa(offset)},
		"limit":   {strconv.Itoa(limit)},
	}
	if event != "" {
		q.Set("event", event)
	}
	if len(transport) > 0 {
		q.Set("transport_types", strings.Join(transport, ","))
	}

	var resp ScheduleResponse
	return &resp, c.do(ctx, "schedule/", q, &resp)
}
