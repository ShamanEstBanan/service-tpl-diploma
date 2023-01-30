package usecase

import (
	"context"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
	"strconv"
)

func (s *service) LoadOrder(ctx context.Context, orderID int, userID string) error {
	validNum := LunaValidation(orderID)
	if !validNum {
		s.lg.Error(errs.ErrInvalidOrderID.Error())
		return errs.ErrInvalidOrderID
	}
	oi := int64(orderID)
	strOrderID := strconv.FormatInt(oi, 10)
	err := s.storage.LoadOrder(ctx, strOrderID, userID)
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

func (s *service) GetUserOrders(ctx context.Context, userID string) (orders []domain.ResponseOrder, err error) {
	orders, err = s.storage.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}
