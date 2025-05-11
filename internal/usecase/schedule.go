package usecase

import (
	"context"
	"fmt"
	"time"
	"train-backend/internal/infrastructure/cache"

	"train-backend/internal/infrastructure/yandex"
)

type ScheduleUsecase struct {
	cli      *yandex.Client
	cacheTTL time.Duration
}

func NewScheduleUsecase(cli *yandex.Client, ttl time.Duration) *ScheduleUsecase {
	return &ScheduleUsecase{cli: cli, cacheTTL: ttl}
}

func (uc *ScheduleUsecase) Search(ctx context.Context, from, to, date string,
	transport []string, transfers bool, offset, limit int) (*yandex.SearchResponse, error) {

	// Кэшируем по ключу «search:from:to:date:offset:limit»
	if v, ok := cache.Get(ctx, uc.cacheKey("search", from, to, date, offset, limit)); ok {
		return v.(*yandex.SearchResponse), nil
	}

	resp, err := uc.cli.Search(ctx, from, to, date, transport, transfers, offset, limit)
	if err == nil {
		cache.Set(ctx, uc.cacheKey("search", from, to, date, offset, limit), resp, uc.cacheTTL)
	}
	return resp, err
}

func (uc *ScheduleUsecase) Station(ctx context.Context, station, date, event string,
	transport []string, offset, limit int) (*yandex.ScheduleResponse, error) {

	if v, ok := cache.Get(ctx, uc.cacheKey("station", station, date, offset, limit)); ok {
		return v.(*yandex.ScheduleResponse), nil
	}

	resp, err := uc.cli.ScheduleOnStation(ctx, station, date, event, transport, offset, limit)
	if err == nil {
		cache.Set(ctx, uc.cacheKey("station", station, date, offset, limit), resp, uc.cacheTTL)
	}
	return resp, err
}

func (uc *ScheduleUsecase) cacheKey(parts ...any) string {
	return fmt.Sprint(parts...)
}
