package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"train-backend/internal/usecase"
)

type DirectionHandler struct {
	uc *usecase.DirectionUsecase
}

func NewDirectionHandler(r fiber.Router, uc *usecase.DirectionUsecase) {
	h := &DirectionHandler{uc}
	r.Get("/directions", h.List) // GET /api/v1/directions
}

// List godoc – соответствует спецификации paths./directions.get :contentReference[oaicite:4]{index=4}&#8203;:contentReference[oaicite:5]{index=5}
func (h *DirectionHandler) List(c *fiber.Ctx) error {
	off, _ := strconv.Atoi(c.Query("offset", "0"))
	lim, _ := strconv.Atoi(c.Query("limit", "100"))

	dirs, total, err := h.uc.List(c.Context(), off, lim)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"pagination": fiber.Map{"total": total, "offset": off, "limit": lim},
		"data":       dirs,
	})
}
