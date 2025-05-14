package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHealthHandler(r fiber.Router, db *pgxpool.Pool) {
	r.Get("/health", func(c *fiber.Ctx) error {
		if err := db.Ping(c.Context()); err != nil {
			return c.Status(503).JSON(fiber.Map{"db": "down"})
		}
		//if cache.StatsCurrent().Items == 0 {

		//	return c.Status(503).JSON(fiber.Map{"cache": "down"})
		//}
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
