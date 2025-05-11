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

	dirRepo := postgres.NewDirectionRepo(pool)
	dirUC := usecase.NewDirectionUsecase(dirRepo, 3*time.Second)

	app := fiber.New()
	app.Use(middleware.Logger())
	_ = godotenv.Load(".env")
	apiKey := os.Getenv("YANDEX_API_KEY")
	yaCli := yandex.New(apiKey)
	schUC := usecase.NewScheduleUsecase(yaCli, 5*time.Minute)
	stRepo := postgres.NewStationRepo(pool)
	stUC := usecase.NewStation(stRepo, 3*time.Second)
	api := app.Group("/api/v1")
	handler.NewDirectionHandler(api, dirUC)
	handler.NewScheduleHandler(api, schUC)
	handler.NewStationHandler(api, stUC)
	log.Fatal(app.Listen(":8080"))
}
