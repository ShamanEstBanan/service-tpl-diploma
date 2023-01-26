package storage

import (
	"context"
	"fmt"
	"time"
)

func (s *storage) GetAccountWithdrawnPoints(ctx context.Context, accountId string) (withdrawnPoints float32, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT sum(points) FROM withdrawals WHERE account_id = '%s'", accountId)

	err = s.db.QueryRow(ctxT, query).Scan(&withdrawnPoints)
	if err != nil {
		return 0, err
	}
	return withdrawnPoints, nil
}
