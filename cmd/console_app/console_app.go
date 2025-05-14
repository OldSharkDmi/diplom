package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const apiBase = "http://localhost:8080/api/v1"

var moscowLoc *time.Location

func init() {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		fmt.Fprintln(os.Stderr, "warning: cannot load Europe/Moscow, using local:", err)
		loc = time.Local
	}
	moscowLoc = loc
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	from := chooseStation("отправления", reader)
	to := chooseStation("прибытия", reader)

	fmt.Print("Введите дату (YYYY-MM-DD): ")
	dateStr, _ := reader.ReadString('\n')
	dateStr = strings.TrimSpace(dateStr)

	// 1) Прямые
	directs, err := fetchDirect(from.Code, to.Code, dateStr)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error fetching direct:", err)
		return
	}

	if len(directs) > 0 {
		printDirects(from, to, directs, dateStr)
	} else {
		// 2) Маршруты с пересадками
		routes, err := fetchRoutes(from.Code, to.Code, dateStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error fetching routes:", err)
			return
		}
		if len(routes) == 0 {
			fmt.Println("Результатов не найдено.")
			return
		}
		printRoutes(from, to, routes, dateStr)
	}
}

// --- Direct ---

type Direct struct {
	UID   string `json:"train_uid"`
	Title string `json:"thread.title"`
	Dep   string `json:"departure"`
	Arr   string `json:"arrival"`
}

func fetchDirect(from, to, date string) ([]Direct, error) {
	u := fmt.Sprintf("%s/search?from=%s&to=%s&date=%s&limit=100",
		apiBase, url.QueryEscape(from), url.QueryEscape(to), date,
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var raw struct {
		Segments []struct {
			Thread struct {
				UID   string `json:"uid"`
				Title string `json:"title"`
			} `json:"thread"`
			Departure string `json:"departure"`
			Arrival   string `json:"arrival"`
		} `json:"segments"`
	}
	if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
		return nil, err
	}

	out := make([]Direct, len(raw.Segments))
	for i, s := range raw.Segments {
		out[i] = Direct{
			UID:   s.Thread.UID,
			Title: s.Thread.Title,
			Dep:   s.Departure,
			Arr:   s.Arrival,
		}
	}
	return out, nil
}

func printDirects(from, to Station, ds []Direct, date string) {
	if len(ds) > 10 {
		ds = ds[:10]
	}
	fmt.Println("\nПрямые поезда:")
	for i, d := range ds {
		fmt.Printf("[%d] %s\n", i, d.Title)
		fmt.Printf("     %s  %s →  %s в %s  (UID: %s)\n",
			from.Title, formatMSK(d.Dep),
			to.Title, formatMSK(d.Arr),
			d.UID,
		)
	}
	uids := make([]string, len(ds))
	for i, d := range ds {
		uids[i] = d.UID
	}
	selectAndShow(uids, date)
}

// --- Routes with one transfer ---

type RouteSeg struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Dep      string `json:"dep"`
	Arr      string `json:"arr"`
	TrainUID string `json:"train_uid"`
}

type Route struct {
	Segments []RouteSeg `json:"segments"`
	Duration float64    `json:"duration"`
}

func fetchRoutes(from, to, date string) ([]Route, error) {
	u := fmt.Sprintf("%s/routes?from=%s&to=%s&date=%s",
		apiBase, url.QueryEscape(from), url.QueryEscape(to), date,
	)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var out struct {
		Segments []RouteSeg `json:"segments"`
		Duration float64    `json:"duration"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}
	if len(out.Segments) == 0 {
		return nil, nil
	}
	return []Route{{Segments: out.Segments, Duration: out.Duration}}, nil
}

func printRoutes(from, to Station, rs []Route, date string) {
	fmt.Println("\nМаршрут с одной пересадкой:")
	r := rs[0]
	for i, s := range r.Segments {
		legFrom := from.Title
		if i == 1 {
			legFrom = r.Segments[0].To
		}
		legTo := to.Title
		fmt.Printf("  [%d] %s → %s  (UID: %s)\n", i+1,
			formatMSK(s.Dep), formatMSK(s.Arr),
			s.TrainUID,
		)
		fmt.Printf("      %s в %s → %s в %s\n",
			legFrom, formatMSK(s.Dep),
			legTo, formatMSK(s.Arr),
		)
		if i+1 < len(r.Segments) {
			next := parseTime(r.Segments[i+1].Dep)
			cur := parseTime(s.Arr)
			fmt.Printf("      Пересадка: %d мин\n", int(next.Sub(cur).Minutes()))
		}
	}
	uids := []string{r.Segments[0].TrainUID, r.Segments[1].TrainUID}
	selectAndShow(uids, date)
}

// --- common select & detail ---

func selectAndShow(uids []string, date string) {
	fmt.Print("\nВыберите индекс поезда: ")
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	idx, _ := strconv.Atoi(strings.TrimSpace(line))
	if idx < 0 || idx >= len(uids) {
		fmt.Println("Неверный индекс.")
		return
	}
	uid := uids[idx]

	// Status
	st, err := fetchStatus(uid)
	if err != nil {
		fmt.Fprintln(os.Stderr, "status error:", err)
	} else {
		fmt.Printf("\nСтатус поезда %s:\n", uid)
		fmt.Printf("  Отправление: %s\n", formatMSK(st.Departure))
		fmt.Printf("  Прибытие:    %s\n", formatMSK(st.Arrival))
		fmt.Printf("  Задержка:    %d мин\n", st.DelayMin)
		fmt.Printf("  Занятость:   %s\n", st.Occupancy)
	}

	// Stops
	fmt.Println("\nОстановки:")
	stops, err := fetchStops(uid, date)
	if err != nil {
		fmt.Fprintln(os.Stderr, "stops error:", err)
		return
	}
	for _, s := range stops {
		fmt.Printf("%-25s %s → %s\n",
			s.Station.Title,
			formatMSKFull(s.Dep, date),
			formatMSKFull(s.Arr, date),
		)
	}
}

// --- helpers ---

type Station struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

func chooseStation(kind string, reader *bufio.Reader) Station {
	fmt.Printf("Введите часть названия станции %s: ", kind)
	q, _ := reader.ReadString('\n')
	q = strings.TrimSpace(q)
	u := fmt.Sprintf("%s/stations?search=%s&limit=10",
		apiBase, url.QueryEscape(q),
	)
	res, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	var out struct {
		Data []Station `json:"data"`
	}
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		panic(err)
	}
	for i, st := range out.Data {
		fmt.Printf("[%d] %s (%s)\n", i, st.Title, st.Code)
	}
	fmt.Print("Выберите индекс станции: ")
	line, _ := reader.ReadString('\n')
	idx, _ := strconv.Atoi(strings.TrimSpace(line))
	return out.Data[idx]
}

type Status struct {
	Departure string `json:"departure"`
	Arrival   string `json:"arrival"`
	DelayMin  int    `json:"delay_min"`
	Occupancy string `json:"occupancy"`
}

func fetchStatus(uid string) (Status, error) {
	var st Status
	u := apiBase + "/trains/" + uid
	res, err := http.Get(u)
	if err != nil {
		return st, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		body, _ := io.ReadAll(res.Body)
		return st, fmt.Errorf("status %d: %s", res.StatusCode, body)
	}
	err = json.NewDecoder(res.Body).Decode(&st)
	return st, err
}

type Stop struct {
	Station struct{ Title string } `json:"station"`
	Dep     string                 `json:"departure"`
	Arr     string                 `json:"arrival"`
}

func fetchStops(uid, date string) ([]Stop, error) {
	u := fmt.Sprintf("%s/trains/%s/stops?date=%s", apiBase, uid, date)
	res, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var stops []Stop
	err = json.NewDecoder(res.Body).Decode(&stops)
	return stops, err
}

func formatMSK(ts string) string {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	return t.In(moscowLoc).Format("15:04")
}
func formatMSKFull(ts, baseDate string) string {
	if ts == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return ts
	}
	base, _ := time.ParseInLocation("2006-01-02", baseDate, moscowLoc)
	if t.Before(base) {
		t = t.Add(24 * time.Hour)
	}
	return t.In(moscowLoc).Format("2006-01-02 15:04")
}
func parseTime(ts string) time.Time {
	t, _ := time.Parse(time.RFC3339, ts)
	return t.In(moscowLoc)
}
