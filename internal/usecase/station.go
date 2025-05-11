package usecase

import (
	"context"
	"time"

	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type Station struct {
	repo           repository.StationRepository
	contextTimeout time.Duration
}

func NewStation(r repository.StationRepository, t time.Duration) *Station {
	return &Station{repo: r, contextTimeout: t}
}

func (uc *Station) Search(ctx context.Context, q string, limit int) ([]model.Station, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	return uc.repo.Search(ctx, q, limit)
}
