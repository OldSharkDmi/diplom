package handler

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"train-backend/internal/usecase"
)

type RouteHandler struct{ uc *usecase.RouteUsecase }

func NewRouteHandler(r fiber.Router, uc *usecase.RouteUsecase) {
	h := &RouteHandler{uc}
	r.Get("/routes", h.Find)
}

func (h *RouteHandler) Find(c fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	date, _ := time.Parse("2006-01-02", c.Query("date"))
	rt, err := h.uc.Find(c.Context(), from, to, date)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(rt)
}
