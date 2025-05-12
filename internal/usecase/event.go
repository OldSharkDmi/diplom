package usecase

import (
	"context"
	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type EventUsecase struct{ repo repository.EventRepository }

func NewEventUsecase(r repository.EventRepository) *EventUsecase { return &EventUsecase{r} }

func (uc *EventUsecase) Store(ctx context.Context, e *model.Event) (*model.Event, error) {
	return e, uc.repo.Store(ctx, e)
}
