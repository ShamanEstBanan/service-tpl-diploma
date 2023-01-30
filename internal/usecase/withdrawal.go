package usecase

import (
	"context"
	"strconv"

	"go.uber.org/zap"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
)

func (s *service) MakeWithdrawn(ctx context.Context, userID string, orderID string, sum float32) (err error) {
	orderINT, err := strconv.Atoi(orderID)
	if err != nil {
		s.lg.Error("Error converting orderID to int", zap.Error(err))
		return err
	}
	validNum := LunaValidation(orderINT)
	if !validNum {
		s.lg.Error(errs.ErrInvalidOrderID.Error())
		return errs.ErrInvalidOrderID
	}

	err = s.storage.MakeWithdrawn(ctx, userID, orderID, sum)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetUserWithdrawals(ctx context.Context, userID string) ([]domain.Withdrawal, error) {
	withdrawals, err := s.storage.GetUserWithdrawals(ctx, userID)
	if err != nil {
		return nil, err
	}
	return withdrawals, nil
}
