package middleware

import (
	"github.com/gofiber/fiber/v3"
	"log"
)

func Logger() fiber.Handler {
	return func(c fiber.Ctx) error {
		err := c.Next()
		log.Printf("%s %s %d", c.Method(), c.Path(), c.Response().StatusCode())
		return err
	}
}
