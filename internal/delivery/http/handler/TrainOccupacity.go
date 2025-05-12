package handler

import (
	"github.com/gofiber/fiber/v3"
	"train-backend/internal/usecase"
)

func TrainOccupancy(uc *usecase.Train) fiber.Handler {
	return func(c fiber.Ctx) error {
		o, err := uc.Occupancy(c.Context(), c.Params("uid"))
		if err != nil {
			return fiber.NewError(404, err.Error())
		}
		return c.JSON(o)
	}
}
