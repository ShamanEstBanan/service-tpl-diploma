package storage

import (
	"context"
	"fmt"
	"time"
)

func (s *storage) GetAccountBalance(ctx context.Context, accountId string) (balance float32, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT balance FROM accounts WHERE id = '%s'", accountId)

	err = s.db.QueryRow(ctxT, query).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
