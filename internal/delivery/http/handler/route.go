package handler

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"

	"train-backend/internal/usecase"
)

type RouteHandler struct{ uc *usecase.RouteUsecase }

func NewRouteHandler(r fiber.Router, uc *usecase.RouteUsecase) {
	h := &RouteHandler{uc}
	r.Get("/routes", h.Find)
	r.Get("/routes/:id", h.Get)
}

func (h *RouteHandler) Find(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")
	date, _ := time.Parse("2006-01-02", c.Query("date"))
	rt, err := h.uc.Find(c.Context(), from, to, date)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(rt)
}
func (h *RouteHandler) Get(c *fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	rt, err := h.uc.ByID(c.Context(), id) // используем тот же use-case
	if err != nil {
		return fiber.NewError(404, err.Error())
	}
	return c.JSON(rt)
}
