package usecase

import (
	"context"

	"service-tpl-diploma/internal/domain"
)

func (s *service) GetUserBalance(ctx context.Context, accountID string) (balance domain.UserBalanceResponse, err error) {
	balance = domain.UserBalanceResponse{}
	current, err := s.storage.GetAccountBalance(ctx, accountID)
	if err != nil {
		return balance, err
	}
	balance.Current = current
	points, err := s.storage.GetAccountWithdrawnPoints(ctx, accountID)
	if err != nil {
		return balance, err
	}
	balance.Withdrawn = points
	return balance, nil
}
