package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
)

func (s *storage) GetAccountWithdrawnPoints(ctx context.Context, accountID string) (points float32, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var wp float32
	query := fmt.Sprintf("SELECT sum(points) FROM withdrawals WHERE account_id = '%s'", accountID)
	err = s.db.QueryRow(ctxT, query).Scan(&wp)

	if errors.Is(err, errs.ErrNoPoints) {
		return 0, errs.ErrNoPoints
	}

	return wp, nil
}

func (s *storage) MakeWithdrawn(ctx context.Context, accountID string, orderID string, amount float32) (err error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	tx, err := s.db.Begin(ctxT)
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		} else {
			tx.Commit(context.TODO())
		}
	}()

	var balance float32
	qTakeBalance := fmt.Sprintf("SELECT balance FROM accounts WHERE id = '%s'", accountID)
	err = tx.QueryRow(ctxT, qTakeBalance).Scan(&balance)
	if err != nil {
		return err
	}
	newBalance := balance - amount
	if newBalance < 0 {
		return errs.ErrNotEnoughPoints
	}

	qUpdBalance := fmt.Sprintf("UPDATE accounts SET balance=%v WHERE id = '%s'", newBalance, accountID)
	_, err = tx.Exec(ctxT, qUpdBalance)
	if err != nil {
		return err
	}

	qSetWithDrawn := "INSERT INTO withdrawals (account_id, order_id, points) VALUES($1,$2,$3)"
	_, err = tx.Exec(ctxT, qSetWithDrawn, accountID, orderID, amount)
	switch e := err.(type) {
	case *pgconn.PgError:
		switch e.Code {
		case pgerrcode.UniqueViolation:
			return errs.ErrWithdrawnAlreadyDoneForThisOrder
		}
	}

	return nil
}

func (s *storage) GetUserWithdrawals(ctx context.Context, userID string) ([]domain.Withdrawal, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(
		"SELECT order_id, points, updated_at FROM withdrawals where account_id = '%s' ORDER BY updated_at DESC",
		userID)
	rows, err := s.db.Query(ctxT, query)
	if err != nil {
		s.lg.Error("ERROR db taking user's withdrawals:", zap.Error(err))
		return nil, err
	}
	defer rows.Close()
	var withdrawals []domain.Withdrawal
	for rows.Next() {
		var processedAt time.Time
		w := domain.Withdrawal{}
		err = rows.Scan(&w.Order, &w.Sum, &processedAt)
		if err != nil {
			s.lg.Error("ERROR scan withdrawals:", zap.Error(err))
			return nil, err
		}
		w.ProcessedAt = processedAt.Format(time.RFC3339)
		withdrawals = append(withdrawals, w)
	}
	return withdrawals, nil
}
