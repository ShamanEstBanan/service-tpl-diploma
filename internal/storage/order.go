package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
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

type UserOrder struct {
	OrderId    string `db:"id"`
	Status     string `db:"status"`
	Accrual    int64  `db:"accrual"`
	UploadedAt string `db:"uploaded_at"`
}

func (s *storage) GetUserOrders(ctx context.Context, userID string) (orders []domain.ResponseOrder, err error) {
	ctxCancel, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := fmt.Sprintf("SELECT * FROM orders WHERE account_id='%s' ORDER BY uploaded_at DESC", userID)
	rows, err := s.db.Query(ctxCancel, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order domain.ResponseOrder
		var accId string
		var uploadedAt, updatedAt time.Time
		err = rows.Scan(&order.Number, &accId, &order.Status, &order.Accrual, &uploadedAt, &updatedAt)
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
