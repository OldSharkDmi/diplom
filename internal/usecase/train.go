package usecase

import (
	"context"
	"time"

	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
	"train-backend/internal/infrastructure/yandex"
)

type Train struct {
	repo repository.TrainRepository
	ya   *yandex.Client
	ttl  time.Duration
}

func NewTrain(r repository.TrainRepository, ya *yandex.Client, ttl time.Duration) *Train {
	return &Train{r, ya, ttl}
}
func (uc *Train) Occupancy(ctx context.Context, uid string) (*model.Occupancy, error) {
	return uc.repo.Occupancy(ctx, uid)
}
func (uc *Train) Status(ctx context.Context, uid string) (*model.TrainStatus, error) {
	// from DB
	if st, _ := uc.repo.Get(ctx, uid); st != nil {
		return st, nil
	}
	// fallback — запрос к Яндексу (нет публичного энд-пойнта, делаем search by uid)
	resp, err := uc.ya.Search(ctx, "", "", time.Now().Format("2006-01-02"),
		[]string{"suburban"}, false, 0, 1)
	if err != nil || len(resp.Segments) == 0 {
		return nil, err
	}
	seg := resp.Segments[0]
	st := &model.TrainStatus{
		UID:        uid,
		Departure:  seg.Departure,
		Arrival:    seg.Arrival,
		UpdatedISO: time.Now().Format(time.RFC3339),
		DelayMin:   0,
		Occupancy:  "unknown",
	}
	_ = uc.repo.Save(ctx, st)
	return st, nil
}
