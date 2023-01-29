package usecase

import (
	"context"
	"service-tpl-diploma/internal/domain"

	"go.uber.org/zap"
)

type storage interface {
	CreateUser(ctx context.Context, user domain.NewUser) (string, error)
	CheckUser(ctx context.Context, user domain.AuthUser) (exist string, err error)
	LoadOrder(ctx context.Context, orderID string, userID string) (err error)
	GetUserOrders(ctx context.Context, userID string) (orders []domain.ResponseOrder, err error)
	GetAccountBalance(ctx context.Context, accountId string) (balance float32, err error)
	GetAccountWithdrawnPoints(ctx context.Context, accountId string) (withdrawnPoints float32, err error)
	MakeWithdrawn(ctx context.Context, userID string, orderID string, amount float32) (err error)
	GetUserWithdrawals(ctx context.Context, userID string) ([]domain.Withdrawal, error)
}

type service struct {
	lg      *zap.Logger
	storage storage
	jobs    chan domain.Job
}

func New(logger *zap.Logger, storage storage) *service {
	return &service{
		lg:      logger,
		storage: storage,
	}
}
