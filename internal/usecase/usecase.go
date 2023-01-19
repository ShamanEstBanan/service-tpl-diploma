package usecase

import (
	"context"
	"service-tpl-diploma/internal/app/domain"

	"go.uber.org/zap"
)

type storage interface {
	CreateUser(ctx context.Context, user domain.NewUser) error
	CheckUser(ctx context.Context, user domain.AuthUser) (exist string, err error)
}

type service struct {
	lg      *zap.Logger
	storage storage
}

func New(logger *zap.Logger, storage storage) *service {
	return &service{
		lg:      logger,
		storage: storage,
	}
}
