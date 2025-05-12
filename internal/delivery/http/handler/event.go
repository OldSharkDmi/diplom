package handler

import (
	"github.com/gofiber/fiber/v3"
	"train-backend/internal/domain/model"
	"train-backend/internal/usecase"
)

type EventHandler struct{ uc *usecase.EventUsecase }

func NewEventHandler(r fiber.Router, uc *usecase.EventUsecase) {
	h := &EventHandler{uc}
	r.Post("/events", h.Store)
}

func (h *EventHandler) Store(c fiber.Ctx) error {
	var req struct {
		DeviceID string      `json:"device_id,omitempty"`
		Type     string      `json:"event_type"`
		Payload  interface{} `json:"payload"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return fiber.ErrBadRequest
	}
	ev := &model.Event{
		DeviceID: req.DeviceID,
		Type:     req.Type,
		Payload:  req.Payload,
	}
	ev, err := h.uc.Store(c.Context(), ev)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.Status(201).JSON(ev)
}
