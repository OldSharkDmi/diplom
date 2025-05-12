package usecase

import (
	"context"
	"time"

	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
	"train-backend/internal/infrastructure/yandex"
)

type RouteUsecase struct {
	repo    repository.RouteRepository
	ya      *yandex.Client
	timeout time.Duration
}

func NewRouteUsecase(r repository.RouteRepository, ya *yandex.Client, t time.Duration) *RouteUsecase {
	return &RouteUsecase{repo: r, ya: ya, timeout: t}
}

func (uc *RouteUsecase) Find(ctx context.Context, from, to string, date time.Time) (*model.Route, error) {
	// 1. cache
	if rt, _ := uc.repo.GetCache(ctx, from, to, date); rt != nil {
		return rt, nil
	}

	// 2. brute DFS 1-пересадка (пример)
	types := []string{"suburban"}
	srch, err := uc.ya.Search(ctx, from, to, date.Format("2006-01-02"), types, false, 0, 50)
	if err == nil && len(srch.Segments) > 0 {
		rt := convertDirect(srch)
		_ = uc.repo.SaveCache(ctx, from, to, date, rt)
		return rt, nil
	}

	// 3. пересадка: перебираем станции-кандидаты
	sr1, _ := uc.ya.Search(ctx, from, "", date.Format("2006-01-02"), types, false, 0, 100)
	best := &model.Route{Duration: 1<<31 - 1}
	for _, seg1 := range sr1.Segments {
		mid := seg1.To.Code
		sr2, err := uc.ya.Search(ctx, mid, to, date.Format("2006-01-02"), types, false, 0, 50)
		if err != nil || len(sr2.Segments) == 0 {
			continue
		}
		seg2 := sr2.Segments[0]
		dur := seg1.Duration + seg2.Duration
		if dur < best.Duration {
			best = &model.Route{
				Segments: []model.RouteSegment{
					{FromCode: from, ToCode: mid, Dep: seg1.Departure, Arr: seg1.Arrival, TrainUID: seg1.Thread.UID},
					{FromCode: mid, ToCode: to, Dep: seg2.Departure, Arr: seg2.Arrival, TrainUID: seg2.Thread.UID},
				},
				Duration: dur,
			}
		}
	}
	_ = uc.repo.SaveCache(ctx, from, to, date, best)
	return best, nil
}

func convertDirect(s *yandex.SearchResponse) *model.Route {
	seg := s.Segments[0]
	return &model.Route{
		Segments: []model.RouteSegment{
			{FromCode: seg.From.Code, ToCode: seg.To.Code, Dep: seg.Departure,
				Arr: seg.Arrival, TrainUID: seg.Thread.UID},
		},
		Duration: seg.Duration,
	}
}
