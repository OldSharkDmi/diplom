package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"train-backend/internal/usecase"
)

type SubHandler struct{ uc *usecase.SubscriptionUsecase }

func NewSubHandler(r fiber.Router, uc *usecase.SubscriptionUsecase) {
	h := &SubHandler{uc}
	r.Post("/subscriptions", h.Create)
	r.Delete("/subscriptions/:id", h.Delete)
}

func (h *SubHandler) Create(c fiber.Ctx) error {
	var req struct {
		DeviceToken string `json:"device_token"`
		TrainUID    string `json:"train_uid"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}
	s, err := h.uc.Subscribe(c.Context(), req.DeviceToken, req.TrainUID)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.Status(201).JSON(s)
}

func (h *SubHandler) Delete(c fiber.Ctx) error {
	id, _ := strconv.ParseInt(c.Params("id"), 10, 64)
	if err := h.uc.Unsubscribe(c.Context(), id); err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.SendStatus(204)
}
