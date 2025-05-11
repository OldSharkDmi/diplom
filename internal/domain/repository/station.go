package repository

import (
	"context"
	"train-backend/internal/domain/model"
)

type StationRepository interface {
	Search(ctx context.Context, q string, limit int) ([]model.Station, error)
}
