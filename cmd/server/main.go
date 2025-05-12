package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"train-backend/internal/delivery/http/handler"
	"train-backend/internal/delivery/http/middleware"
	"train-backend/internal/infrastructure/database"
	"train-backend/internal/infrastructure/repository/postgres"
	"train-backend/internal/infrastructure/yandex"
	"train-backend/internal/usecase"
)

func main() {
	ctx := context.Background()
	pool, err := database.NewPool(ctx)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	log.Println("✅ Подключение к базе данных успешно выполнено")

	dirRepo := postgres.NewDirectionRepo(pool)
	dirUC := usecase.NewDirectionUsecase(dirRepo, 3*time.Second)

	app := fiber.New()
	app.Use(middleware.Logger())
	_ = godotenv.Load(".env")
	apiKey := os.Getenv("YANDEX_API_KEY")
	yaCli := yandex.New(apiKey)
	schUC := usecase.NewScheduleUsecase(yaCli, 5*time.Minute)
	routeRepo := postgres.NewRouteRepo(pool)
	routeUC := usecase.NewRouteUsecase(routeRepo, yaCli, 10*time.Second)

	stRepo := postgres.NewStationRepo(pool)
	stUC := usecase.NewStation(stRepo, 3*time.Second)
	trainRepo := postgres.NewTrainRepo(pool)
	trainUC := usecase.NewTrain(trainRepo, yaCli, 5*time.Minute)
	subRepo := postgres.NewSubRepo(pool)
	subUC := usecase.NewSubscriptionUsecase(subRepo)
	eventRepo := postgres.NewEventRepo(pool)
	eventUC := usecase.NewEventUsecase(eventRepo)
	api := app.Group("/api/v1")
	handler.NewDirectionHandler(api, dirUC)
	handler.NewScheduleHandler(api, schUC)
	handler.NewStationHandler(api, stUC)
	handler.NewRouteHandler(api, routeUC)
	handler.NewSubHandler(api, subUC)
	handler.NewHealthHandler(api, pool)
	handler.NewTrainHandler(api, trainUC)
	handler.NewEventHandler(api, eventUC)
	log.Fatal(app.Listen(":8080"))

	go func() { // listener push → WS / FCM
		conn, _ := pool.Acquire(context.Background())
		defer conn.Release()
		_, _ = conn.Exec(context.Background(), `LISTEN train_updates`)
		for {
			if n, err := conn.Conn().WaitForNotification(context.Background()); err == nil {
				uid := n.Payload
				subs, _ := subRepo.ByTrain(context.Background(), uid)
				for _, sub := range subs {
					// TODO: отправить push (FCM / APNS)  – здесь только лог
					log.Printf("[PUSH] train %s → token %s", uid, sub.DeviceToken)
				}
			}
		}
	}()

	go func() { // периодический fetch статусов
		ticker := time.NewTicker(2 * time.Minute)
		for range ticker.C {
			uids, _ := trainRepo.WithSubs(context.Background())
			for _, uid := range uids {
				_, _ = trainUC.Status(context.Background(), uid) // обновит таблицу
			}
		}
	}()
}
