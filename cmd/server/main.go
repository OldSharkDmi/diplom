package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v3"
	"train-backend/internal/delivery/http/handler"
	"train-backend/internal/delivery/http/middleware"
	"train-backend/internal/infrastructure/database"
	"train-backend/internal/infrastructure/repository/postgres"
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

	api := app.Group("/api/v1")
	handler.NewDirectionHandler(api, dirUC)

	log.Fatal(app.Listen(":8080"))
}
