package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"train-backend/internal/infrastructure/yandex"
)

func main() {
	_ = godotenv.Load(".env") // для локального запуска
	ctx := context.Background()

	cli := yandex.New(os.Getenv("YANDEX_API_KEY"))
	stations, err := cli.StationsList(ctx, []string{"suburban"})
	if err != nil {
		log.Fatal(err)
	}

	// собираем DSN из тех же переменных, что использует сервер
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("PGUSER"), os.Getenv("PGPASS"),
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGDB"))

	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	batch := &pgx.Batch{}
	for _, s := range stations {
		batch.Queue(`INSERT INTO stations
		(code,title,station_type,transport_type,latitude,longitude,settlement_code)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (code) DO UPDATE SET title=EXCLUDED.title`,
			s.Code, s.Title, s.Type, s.Transport,
			s.Latitude, s.Longitude, s.SettlementCode)
	}
	if err := conn.SendBatch(ctx, batch).Close(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(os.Stdout).Encode(map[string]int{"imported": len(stations)})
}
