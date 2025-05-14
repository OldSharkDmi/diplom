package handler

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"train-backend/internal/usecase"
)

// queryBool парсит булевый параметр с дефолтным значением.
func queryBool(c *fiber.Ctx, key string, def bool) bool {
	val := c.Query(key, "")
	if val == "" {
		return def
	}
	b, _ := strconv.ParseBool(val)
	return b
}

type ScheduleHandler struct {
	uc *usecase.ScheduleUsecase
}

// NewScheduleHandler регистрирует два роутинга:
//
//	GET /api/v1/search        — поиск «точка→точка»
//	GET /api/v1/station/:code — расписание на станции
func NewScheduleHandler(r fiber.Router, uc *usecase.ScheduleUsecase) {
	h := &ScheduleHandler{uc}
	r.Get("/search", h.PointToPoint)
	r.Get("/station/:code", h.OnStation)
}

// PointToPoint — прямой поиск маршрута
func (h *ScheduleHandler) PointToPoint(c *fiber.Ctx) error {
	off, _ := strconv.Atoi(c.Query("offset", "0"))
	lim, _ := strconv.Atoi(c.Query("limit", "100"))

	resp, err := h.uc.Search(
		c.Context(),
		c.Query("from"),
		c.Query("to"),
		c.Query("date"),
		strings.Split(c.Query("transport_types", "suburban"), ","),
		queryBool(c, "transfers", false),
		off, lim,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.JSON(resp)
}

// OnStation — расписание на одной станции
func (h *ScheduleHandler) OnStation(c *fiber.Ctx) error {
	off, _ := strconv.Atoi(c.Query("offset", "0"))
	lim, _ := strconv.Atoi(c.Query("limit", "100"))

	resp, err := h.uc.Station(
		c.Context(),
		c.Params("code"),
		c.Query("date"),
		c.Query("event"),
		strings.Split(c.Query("transport_types", "suburban"), ","),
		off, lim,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
	return c.JSON(resp)
}
