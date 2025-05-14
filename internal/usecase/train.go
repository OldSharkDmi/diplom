package usecase

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
	"train-backend/internal/infrastructure/yandex"
)

/* ---------- конструктор ---------- */

type Train struct {
	repo    repository.TrainRepository
	subRepo repository.SubRepository
	ya      *yandex.Client
	ttl     time.Duration
}

func NewTrain(
	trainRepo repository.TrainRepository,
	subRepo repository.SubRepository,
	ya *yandex.Client,
	ttl time.Duration,
) *Train {
	return &Train{repo: trainRepo, subRepo: subRepo, ya: ya, ttl: ttl}
}

/* ---------- ошибки ---------- */

var ErrNotFound = fiber.NewError(fiber.StatusNotFound, "no status")

/* ---------- публичные методы ---------- */

func (uc *Train) Status(ctx context.Context, uid string) (*model.TrainStatus, error) {
	// 1) Смотрим в БД — только если и dep, и arr непустые, возвращаем:
	if st, _ := uc.repo.Get(ctx, uid); st != nil {
		if st.Departure != "" && st.Arrival != "" {
			return st, nil
		}
		// иначе fall-through — идём за свежими данными
	}

	// 2) Дёргаем Yandex /thread
	date := time.Now().Format("2006-01-02")
	tr, err := uc.ya.Thread(ctx, uid, date)
	if err != nil {
		return nil, ErrNotFound
	}

	// 3) Парсим первую и последнюю остановку и форматируем их в RFC-3339
	var dep, arr string
	now := time.Now()
	if len(tr.Stops) > 0 {
		// Yandex возвращает ISO-8601, например "2025-05-14T22:40:00+03:00"
		if t0, err := time.Parse(time.RFC3339, tr.Stops[0].Dep); err == nil {
			dep = t0.Format(time.RFC3339)
		} else {
			dep = tr.Stops[0].Dep // fallback
		}
		if t1, err := time.Parse(time.RFC3339, tr.Stops[len(tr.Stops)-1].Arr); err == nil {
			arr = t1.Format(time.RFC3339)
		} else {
			arr = tr.Stops[len(tr.Stops)-1].Arr
		}
	}

	// 4) Собираем новую модель со всеми полями
	st := &model.TrainStatus{
		UID:        uid,
		Departure:  dep,
		Arrival:    arr,
		UpdatedISO: now.Format(time.RFC3339),
		DelayMin:   0,
		Occupancy:  "unknown",

		// служебные поля
		ReceivedAt: now,
		UpdatedAt:  now,
	}

	// 5) Сохраняем в БД (INSERT … ON CONFLICT UPDATE)
	_ = uc.repo.Save(ctx, st)
	return st, nil
}

// Stops — точки маршрута.
func (uc *Train) Stops(ctx context.Context, uid, date string) ([]yandex.Stop, error) {
	tr, err := uc.ya.Thread(ctx, uid, date)
	if err != nil {
		return nil, err
	}
	return tr.Stops, nil
}

// Occupancy — прогноз.
func (uc *Train) Occupancy(ctx context.Context, uid string) (*model.Occupancy, error) {
	if occ, err := uc.repo.Occupancy(ctx, uid); err == nil {
		return occ, nil
	}

	// очень простая эвристика
	now := time.Now()
	h := now.Hour()
	level := "low"
	switch {
	case h >= 7 && h < 10 || h >= 17 && h < 20:
		level = "high"
	case h >= 5 && h < 7 || h >= 20 && h < 22:
		level = "medium"
	}
	if subs, err := uc.subRepo.ByTrain(ctx, uid); err == nil {
		switch n := len(subs); {
		case n > 30:
			level = "high"
		case n > 10 && level == "low":
			level = "medium"
		}
	}

	return &model.Occupancy{Level: level, UpdatedAt: now}, nil
}
