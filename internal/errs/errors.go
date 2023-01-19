package errs

import (
	"errors"
	"fmt"
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
