package repository

import "context"
import "train-backend/internal/domain/model"

type DirectionRepository interface {
	Fetch(ctx context.Context, offset, limit int) ([]model.Direction, int, error)
}
