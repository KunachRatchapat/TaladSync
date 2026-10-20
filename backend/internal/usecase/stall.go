package usecase

import "context"

type StallRepository interface {
	ListZones(ctx context.Context) ([]string, error)
}

type StallUsecase struct {
	stalls StallRepository
}

func NewStallUsecase(stalls StallRepository) *StallUsecase {
	return &StallUsecase{stalls: stalls}
}

func (u *StallUsecase) ListZones(ctx context.Context) ([]string, error) {
	return u.stalls.ListZones(ctx)
}
