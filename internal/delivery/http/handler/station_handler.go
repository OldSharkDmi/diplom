package handler

import (
	"github.com/gofiber/fiber/v2"
	"strconv"

	"train-backend/internal/usecase"
)

type StationHandler struct{ uc *usecase.Station }

func NewStationHandler(r fiber.Router, uc *usecase.Station) {
	h := &StationHandler{uc}
	r.Get("/stations", h.Search)
}

func (h *StationHandler) Search(c *fiber.Ctx) error { // Изменено на *fiber.Ctx
	q := c.Query("search", "")
	lim, _ := strconv.Atoi(c.Query("limit", "20"))
	data, err := h.uc.Search(c.Context(), q, lim)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(fiber.Map{"data": data})
}
