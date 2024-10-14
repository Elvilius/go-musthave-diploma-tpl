package balances

import (
	"context"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
)

type Storer interface {
	GetBalance(ctx context.Context, userID uint64) (models.Balance, error)
	Withdraw(ctx context.Context, userID uint64, order string, sum float32) error
	GetWithdraws(ctx context.Context, userID uint64) ([]models.GetWithdraw, error)
}

type Service struct {
	store Storer
}

func New(store Storer) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetBalance(ctx context.Context, userID uint64) (models.Balance, error) {
	return s.store.GetBalance(ctx, userID)
}

func (s *Service) Withdraw(ctx context.Context, userID uint64, order string, sum float32) error {
	return s.store.Withdraw(ctx, userID, order, sum)
}

func (s *Service) GetWithdraws(ctx context.Context, userID uint64) ([]models.GetWithdraw, error) {
	return s.store.GetWithdraws(ctx, userID)
}
