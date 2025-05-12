package handler

import (
	"github.com/gofiber/fiber/v3"
)

func errResp(msg string) fiber.Map {
	return fiber.Map{
		"error":   "bad_request",
		"message": msg,
	}
}
