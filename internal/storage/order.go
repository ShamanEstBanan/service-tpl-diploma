package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"log"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
	"time"
)

func (s *storage) LoadOrder(ctx context.Context, orderID string, userID string) error {
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

func (s *storage) UpdateOrder(ctx context.Context, orderID string, status string, accrual decimal.Decimal) error {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	fmt.Printf("orderId :%s\n, status:%s\n, accrual:%v\n", orderID, status, accrual)

	q := fmt.Sprintf(
		"UPDATE orders "+
			"SET status = '%s', accrual = '%v', uploaded_at = now() "+
			"WHERE id LIKE '%s'", status, accrual, orderID)

	//query := "INSERT INTO orders (id, account_id, status, accrual) VALUES($1, $2, $3, $4)"
	_, err := s.db.Exec(ctxCancel, q)
	if err != nil {
		log.Println(err)
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
		"SELECT id FROM orders WHERE status IN ('%s', '%s')",
		domain.OrderInternalStatusNEW, domain.OrderInternalStatusPROCESSING,
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
