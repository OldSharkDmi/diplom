package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type EventRepository interface {
	Store(ctx context.Context, e *model.Event) error
}
