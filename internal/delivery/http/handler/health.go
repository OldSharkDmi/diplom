package handler

import (
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"train-backend/internal/infrastructure/cache"
)

func NewHealthHandler(r fiber.Router, db *pgxpool.Pool) {
	r.Get("/health", func(c fiber.Ctx) error {
		if err := db.Ping(c.Context()); err != nil {
			return c.Status(503).JSON(fiber.Map{"db": "down"})
		}
		if cache.Stat().Items == 0 { // кэш недоступен/пустой
			return c.Status(503).JSON(fiber.Map{"cache": "down"})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
