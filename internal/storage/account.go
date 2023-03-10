package storage

import (
	"context"
	"fmt"
	"time"
)

func (s *storage) GetAccountBalance(ctx context.Context, accountID string) (balance float32, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT balance FROM accounts WHERE id = '%s'", accountID)

	err = s.db.QueryRow(ctxT, query).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (s *storage) UpdateAccountBalance(ctx context.Context, orderID string, accrual float32) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(
		"UPDATE accounts SET balance=balance+%v WHERE id = (SELECT account_id FROM orders WHERE id = '%s')",
		accrual, orderID)
	_, err := s.db.Exec(ctxT, query)
	if err != nil {
		return err
	}
	return nil
}
