package storage

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
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
			if err != nil {
				return err
			}
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

func (s *storage) UpdateOrder(ctx context.Context, orderID string, status string, accrual float32) error {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	q := fmt.Sprintf(
		"UPDATE orders "+
			"SET status = '%s', accrual = %v, uploaded_at = now() "+
			"WHERE id LIKE '%s'", status, accrual, orderID)

	_, err := s.db.Exec(ctxCancel, q)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (s *storage) GetUserOrders(ctx context.Context, userID string) (orders []domain.ResponseOrder, err error) {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := fmt.Sprintf("SELECT * FROM orders WHERE account_id='%s' ORDER BY uploaded_at DESC", userID)
	rows, err := s.db.Query(ctxCancel, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order domain.ResponseOrder
		var accID string
		var uploadedAt, updatedAt time.Time
		err = rows.Scan(&order.Number, &accID, &order.Status, &order.Accrual, &uploadedAt, &updatedAt)
		if err != nil {
			s.lg.Error("ERROR while parse user order: ", zap.Error(err))
			continue
		}
		order.UploadedAt = uploadedAt.Format(time.RFC3339)
		orders = append(orders, order)
	}
	return orders, nil
}

func (s *storage) GetOrdersForProcessing(ctx context.Context) ([]string, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(
		"SELECT id, status FROM orders WHERE status IN ('%s', '%s')",
		domain.OrderInternalStatusNEW, domain.OrderInternalStatusPROCESSING,
	)
	rows, err := s.db.Query(ctxT, query)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var orders []string
	for rows.Next() {
		var orderID, status string
		err = rows.Scan(&orderID, &status)
		if err != nil {
			return nil, err
		}
		if status == domain.OrderInternalStatusNEW {
			queryStatusUpdate := fmt.Sprintf(
				"UPDATE orders SET status = '%s' WHERE id = '%s'",
				domain.OrderInternalStatusPROCESSING, orderID)
			_, err = s.db.Exec(ctxT, queryStatusUpdate)
			if err != nil {
				return nil, err
			}
		}
		orders = append(orders, orderID)
	}
	return orders, nil
}
