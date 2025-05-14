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

	stations, err := uc.repo.Search(ctx, q, limit)
	if err != nil {
		return nil, err
	}

	// Фильтрация только по пригородным поездам
	filtered := make([]model.Station, 0, len(stations))
	for _, s := range stations {
		if s.Transport == "train" {
			filtered = append(filtered, s)
		}
	}
	return filtered, nil
}
