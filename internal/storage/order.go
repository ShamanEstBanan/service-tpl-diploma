package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
	"time"
)

func (s *storage) LoadOrder(ctx context.Context, orderID int, userID string) error {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO orders (id, account_id) VALUES($1, $2)"
	_, err := s.db.Exec(ctxCancel, query, orderID, userID)
	switch e := err.(type) {
	case *pgconn.PgError:
		switch e.Code {
		case pgerrcode.UniqueViolation:
			var isUserOrder bool
			q := fmt.Sprintf("SELECT EXISTS(SELECT * from orders WHERE id='%v' AND account_id='%s')", orderID, userID)
			err = s.db.QueryRow(ctxCancel, q).Scan(&isUserOrder)
			if isUserOrder {
				return errs.ErrOrderAlreadyUploaded
			} else {
				return errs.ErrOrderAlreadyExist
			}
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) UpdateOrder(ctx context.Context, orderID int, userID string, status string, accural decimal.Decimal) error {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := "INSERT INTO orders (id, accountID) VALUES($1, $2)"
	_, err := s.db.Exec(ctxCancel, query, orderID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetUserOrders(ctx context.Context, userID string) error {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := "INSERT INTO orders (id, accountID) VALUES($1, $2)"
	_, err := s.db.Exec(ctxCancel, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetOrdersForProcessing(ctx context.Context) ([]string, error) {
	ctxCancel, cancel := context.WithTimeout(ctx, 50*time.Second)
	defer cancel()
	var orders []string
	var orderId string
	query := fmt.Sprintf(
		"SELECT id FROM orders WHERE status IN ('%s', '%s')", domain.OrderStatusNEW, domain.OrderStatusPROCESSING,
	)
	rows, err := s.db.Query(ctxCancel, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&orderId)
		orders = append(orders, orderId)
	}
	return orders, nil
}
