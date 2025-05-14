package usecase

import (
	"context"
	"errors"
	"math"
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
	// 1. кэш
	if rt, _ := uc.repo.GetCache(ctx, from, to, date); rt != nil {
		return rt, nil
	}

	trySearch := func(limit int, types []string) (*yandex.SearchResponse, error) {
		// отдельный CONTEXT с таймаутом
		cctx, cancel := context.WithTimeout(ctx, uc.timeout)
		defer cancel()
		return uc.ya.Search(cctx, from, to, date.Format("2006-01-02"), types, false, 0, limit)
	}

	/* ---------- прямой поиск ---------- */
	srch, err := trySearch(10, []string{"suburban", "train"})
	if err == nil && len(srch.Segments) > 0 {
		rt := convertDirect(srch)
		_ = uc.repo.SaveCache(ctx, from, to, date, rt)
		return rt, nil
	}
	// Если Яндекс вернул 5xx — отдадим 502 клиенту
	if err != nil {
		return nil, errors.New("yandex upstream error: " + err.Error())
	}

	/* ---------- вариант с одной пересадкой ---------- */
	sr1, _ := trySearch(100, []string{"suburban", "train"})
	best := &model.Route{Duration: math.MaxFloat64}

	for _, seg1 := range sr1.Segments {
		mid := seg1.To.Code
		sr2, _ := uc.ya.Search(ctx, mid, to, date.Format("2006-01-02"),
			[]string{"suburban", "train"}, false, 0, 50)
		if len(sr2.Segments) == 0 {
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

func (uc *RouteUsecase) ByID(ctx context.Context, id int64) (*model.Route, error) {
	return uc.repo.ByID(ctx, id)
}

func convertDirect(s *yandex.SearchResponse) *model.Route {
	seg := s.Segments[0]
	return &model.Route{
		Segments: []model.RouteSegment{
			{FromCode: seg.From.Code, ToCode: seg.To.Code, Dep: seg.Departure, Arr: seg.Arrival, TrainUID: seg.Thread.UID},
		},
		Duration: seg.Duration,
	}
}
