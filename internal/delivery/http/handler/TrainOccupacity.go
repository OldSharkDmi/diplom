package handler

import (
	"github.com/gofiber/fiber/v2"
	"train-backend/internal/usecase"
)

func TrainOccupancy(uc *usecase.Train) fiber.Handler {
	return func(c *fiber.Ctx) error {
		o, err := uc.Occupancy(c.Context(), c.Params("uid"))
		if err != nil {
			return fiber.NewError(500, err.Error())
		}
		if o == nil { // нет данных
			return c.SendStatus(204)
		}
		return c.JSON(o)
	}
}
