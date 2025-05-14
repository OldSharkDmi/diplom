package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"train-backend/internal/usecase"
)

type TrainHandler struct{ uc *usecase.Train }

func NewTrainHandler(r fiber.Router, uc *usecase.Train) {
	h := &TrainHandler{uc}
	r.Get("/trains/:uid/stops", h.Stops)
	r.Get("/trains/:uid/occupancy", h.Occupancy)
	r.Get("/trains/:uid", h.Status)
}

/* ----- handlers ----- */

func (h *TrainHandler) Status(c *fiber.Ctx) error {
	st, err := h.uc.Status(c.Context(), c.Params("uid"))
	return replyJSON(c, st, err)
}

func (h *TrainHandler) Occupancy(c *fiber.Ctx) error {
	occ, err := h.uc.Occupancy(c.Context(), c.Params("uid"))
	return replyJSON(c, occ, err)
}

func (h *TrainHandler) Stops(c *fiber.Ctx) error {
	date := c.Query("date", time.Now().Format("2006-01-02"))
	stops, err := h.uc.Stops(c.Context(), c.Params("uid"), date)
	return replyJSON(c, stops, err)
}

/* ----- helper ----- */

func replyJSON(c *fiber.Ctx, data any, err error) error {
	if err != nil {
		return err // err уже может быть *fiber.Error
	}
	return c.JSON(data)
}
