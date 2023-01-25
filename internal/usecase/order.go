package usecase

import (
	"context"
	"github.com/shopspring/decimal"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
	"strconv"
	"time"
)

func (s *service) LoadOrder(ctx context.Context, orderID int, userId string) error {
	validNum := LunaValidation(orderID)
	if !validNum {
		s.lg.Error(errs.ErrInvalidOrderID.Error())
		return errs.ErrInvalidOrderID
	}
	oi := int64(orderID)
	strOrderID := strconv.FormatInt(oi, 10)
	err := s.storage.LoadOrder(ctx, strOrderID, userId)
	if err != nil {
		s.lg.Error(err.Error())
		return err
	}
	return nil
}

// LunaValidation Valid check number is valid or not based on Luhn algorithm
func LunaValidation(number int) bool {
	return (number%10+checksum(number/10))%10 == 0
}

func checksum(number int) int {
	var luhn int
	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}

func (s *service) GetUserOrders(ctx context.Context, userID string) (orders []domain.Order, err error) {
	ss, _ := decimal.NewFromString("901.1")
	orders = []domain.Order{
		{
			Number:     "1",
			Status:     "1",
			UploadedAt: time.Now(),
		},
		{
			Number:     "2",
			Status:     "3",
			Accrual:    ss,
			UploadedAt: time.Now(),
		},
	}
	return orders, nil
}
