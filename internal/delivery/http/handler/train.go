package handler

import (
	"github.com/gofiber/fiber/v3"
	"train-backend/internal/usecase"
)

type TrainHandler struct{ uc *usecase.Train }

func NewTrainHandler(r fiber.Router, uc *usecase.Train) {
	h := &TrainHandler{uc}
	r.Get("/trains/:uid", h.Status)
}

func (h *TrainHandler) Status(c fiber.Ctx) error {
	uid := c.Params("uid")
	st, err := h.uc.Status(c.Context(), uid)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.JSON(st)
}
