package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"train-backend/internal/infrastructure/yandex"
)

func main() {
	_ = godotenv.Load(".env")
	ctx := context.Background()
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// 1. Получаем станции ----------------------------------------------------
	cli := yandex.New(os.Getenv("YANDEX_API_KEY"))
	cli.SetTimeout(60 * time.Second)

	log.Println("⌛  Запрашиваем станции…")
	t0 := time.Now()
	stations, err := cli.StationsList(ctx, []string{"suburban"})
	if err != nil {
		log.Fatalf("API error: %v", err)
	}
	log.Printf("✅  %d объектов за %v", len(stations), time.Since(t0))

	// 2. Подключаемся к Postgres --------------------------------------------
	dsn := fmt.Sprintf(
		"host=127.0.0.1 port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=5",
		os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASS"), os.Getenv("PGDB"),
	)
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("PG connect: %v", err)
	}
	defer conn.Close(ctx)

	// 3. Чистим таблицу ------------------------------------------------------
	if _, err := conn.Exec(ctx, `TRUNCATE TABLE stations`); err != nil {
		log.Fatalf("TRUNCATE: %v", err)
	}

	// 4. Формируем данные: пропускаем пустой code и дубликаты ----------------
	seen := make(map[string]struct{}, len(stations))
	rows := make([][]any, 0, len(stations))

	for _, s := range stations {
		if s.Code == "" { // пропускаем записи без кода
			continue
		}
		if _, dup := seen[s.Code]; dup { // убираем дубликаты
			continue
		}
		seen[s.Code] = struct{}{}

		lat, _ := s.Latitude.Float64()
		lon, _ := s.Longitude.Float64()
		rows = append(rows, []any{
			s.Code, s.Title, s.Type, s.Transport, lat, lon, s.SettlementCode,
		})
	}
	log.Printf("🚀  COPY %d строк…", len(rows))

	// 5. COPY ---------------------------------------------------------------
	startCopy := time.Now()
	n, err := conn.CopyFrom(
		ctx,
		pgx.Identifier{"stations"},
		[]string{"code", "title", "station_type", "transport_type", "latitude", "longitude", "settlement_code"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		log.Fatalf("CopyFrom: %v", err)
	}
	el := time.Since(startCopy)
	log.Printf("✅  %d строк за %v (%.0f rows/s)", n, el, float64(n)/el.Seconds())

	// 6. Итог ---------------------------------------------------------------
	_ = json.NewEncoder(os.Stdout).Encode(map[string]any{
		"imported": n,
		"elapsed":  time.Since(t0).String(),
	})
}
