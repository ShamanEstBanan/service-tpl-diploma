package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"service-tpl-diploma/internal/errs"
	"time"
)

func (s *storage) GetAccountWithdrawnPoints(ctx context.Context, accountId string) (withdrawnPoints float32, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT sum(points) FROM withdrawals WHERE account_id = '%s'", accountId)
	err = s.db.QueryRow(ctxT, query).Scan(&withdrawnPoints)

	erNull := errors.New("cannot scan NULL into *float32")
	if err != nil {
		if errors.As(err, &erNull) {
			return 0, errs.ErrNoWithdrawn
		}
		return 0, err
	}
	return withdrawnPoints, nil
}

func (s *storage) MakeWithdrawn(ctx context.Context, accountID string, orderID string, amount float32) (err error) {
	//ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	//defer cancel()
	ctxT := context.Background()
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
		return errs.ErrNotEnoughtPoints
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
