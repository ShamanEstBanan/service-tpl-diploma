package usecase

import (
	"context"
	"service-tpl-diploma/internal/domain"

	"go.uber.org/zap"
)

type storage interface {
	CreateUser(ctx context.Context, user domain.NewUser) error
	CheckUser(ctx context.Context, user domain.AuthUser) (exist string, err error)
	LoadOrder(ctx context.Context, orderID string, userID string) (err error)
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
