package errs

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAccessToken               = errors.New("invalid auth token")
	ErrInvalidLoginOrPassword           = errors.New("login or password is invalid")
	ErrLoginIsEmpty                     = errors.New("login is empty")
	ErrLoginAlreadyExist                = errors.New("login already exist")
	ErrPasswordIsEmpty                  = errors.New("password is empty")
	ErrInvalidOrderID                   = errors.New("invalid order number")
	ErrOrderAlreadyExist                = errors.New("orderID already exist")
	ErrOrderAlreadyUploaded             = errors.New("order already uploaded ")
	ErrWithdrawnAlreadyDoneForThisOrder = errors.New("withdrawn already done for this order")
	ErrNotEnoughPoints                  = errors.New("not enough points on balance")
	ErrNoWithdrawn                      = errors.New("no withdrawn")
	ErrNoPoints                         = errors.New("can't scan into dest[0]: cannot scan NULL into *float32")
)

type SQLError struct {
	Code string
	Err  error
}

func (se *SQLError) Error() string {
	return fmt.Sprintf("%v", se.Code)
}

func NewSQLError(code string) error {
	return &SQLError{
		Code: code,
	}
}

func (se *SQLError) Unwrap() error {
	return errors.New(se.Code)
}
