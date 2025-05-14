package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"train-backend/internal/delivery/http/handler"
	"train-backend/internal/delivery/http/middleware"
	"train-backend/internal/delivery/ws"
	"train-backend/internal/domain/repository"
	"train-backend/internal/infrastructure/database"
	"train-backend/internal/infrastructure/repository/postgres"
	"train-backend/internal/infrastructure/yandex"
	"train-backend/internal/usecase"
)

func main() {
	_ = godotenv.Load(".env")
	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	log.Println("✅ подключились к Postgres")

	// строим сервисы и use-cases...
	apiKey := os.Getenv("YANDEX_API_KEY")
	yaCli := yandex.New(apiKey)
	trainRepo := postgres.NewTrainRepo(pool)
	subRepo := postgres.NewSubRepo(pool)
	dirUC := usecase.NewDirectionUsecase(postgres.NewDirectionRepo(pool), 3*time.Second)
	stUC := usecase.NewStation(postgres.NewStationRepo(pool), 3*time.Second)
	routeUC := usecase.NewRouteUsecase(postgres.NewRouteRepo(pool), yaCli, 10*time.Second)
	trainUC := usecase.NewTrain(trainRepo, subRepo, yaCli, 5*time.Minute)
	subUC := usecase.NewSubscriptionUsecase(postgres.NewSubRepo(pool))
	eventUC := usecase.NewEventUsecase(postgres.NewEventRepo(pool))
	schUC := usecase.NewScheduleUsecase(yaCli, 5*time.Minute)

	// HTTP-сервер
	app := fiber.New()
	app.Use(middleware.Logger())

	api := app.Group("/api/v1")
	handler.NewHealthHandler(api, pool)
	handler.NewDirectionHandler(api, dirUC)
	handler.NewStationHandler(api, stUC)
	handler.NewScheduleHandler(api, schUC)
	handler.NewRouteHandler(api, routeUC)
	handler.NewTrainHandler(api, trainUC)
	handler.NewSubHandler(api, subUC)
	handler.NewEventHandler(api, eventUC)

	api.Get("/trains/:uid/occupancy", handler.TrainOccupancy(trainUC))
	app.Get("/ws/push", ws.Handler(postgres.NewTrainRepo(pool)))

	// фоновые задачи — здесь repo, а не usecase!
	go startPushFCM(pool, postgres.NewSubRepo(pool))
	go startStatusFetcher(postgres.NewTrainRepo(pool), trainUC)

	log.Fatal(app.Listen(":8080"))
}

func startPushFCM(pool *pgxpool.Pool, subRepo repository.SubRepository) {
	conn, _ := pool.Acquire(context.Background())
	defer conn.Release()
	_, _ = conn.Exec(context.Background(), `LISTEN train_updates`)
	for {
		if n, err := conn.Conn().WaitForNotification(context.Background()); err == nil {
			uid := n.Payload
			subs, _ := subRepo.ByTrain(context.Background(), uid)
			for _, sub := range subs {
				log.Printf("[PUSH] train %s → token %s", uid, sub.DeviceToken)
			}
		}
	}
}

func startStatusFetcher(trainRepo repository.TrainRepository, trainUC *usecase.Train) {
	ticker := time.NewTicker(2 * time.Minute)
	for range ticker.C {
		uids, _ := trainRepo.WithSubs(context.Background())
		for _, uid := range uids {
			_, _ = trainUC.Status(context.Background(), uid)
		}
	}
}
