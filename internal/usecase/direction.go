package usecase

import (
	"context"
	"time"
	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type DirectionUsecase struct {
	repo           repository.DirectionRepository
	contextTimeout time.Duration
}

func NewDirectionUsecase(r repository.DirectionRepository, timeout time.Duration) *DirectionUsecase {
	return &DirectionUsecase{repo: r, contextTimeout: timeout}
}

func (uc *DirectionUsecase) List(ctx context.Context, offset, limit int) (dirs []model.Direction, total int, err error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	return uc.repo.Fetch(ctx, offset, limit)
}
