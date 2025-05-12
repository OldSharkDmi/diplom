package usecase

import (
	"context"
	"train-backend/internal/domain/model"
	"train-backend/internal/domain/repository"
)

type SubscriptionUsecase struct{ repo repository.SubRepository }

func NewSubscriptionUsecase(r repository.SubRepository) *SubscriptionUsecase {
	return &SubscriptionUsecase{r}
}

func (uc *SubscriptionUsecase) Subscribe(ctx context.Context, token, uid string) (*model.Subscription, error) {
	s := &model.Subscription{DeviceToken: token, TrainUID: uid}
	return s, uc.repo.Create(ctx, s)
}

func (uc *SubscriptionUsecase) Unsubscribe(ctx context.Context, id int64) error {
	return uc.repo.Delete(ctx, id)
}
