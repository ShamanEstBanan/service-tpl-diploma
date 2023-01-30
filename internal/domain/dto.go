package domain

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

const (
	OrderInternalStatusNEW        = "NEW"
	OrderInternalStatusPROCESSING = "PROCESSING"
)

type Job interface {
	Run(ctx context.Context) error
}

type Order struct {
	Number     string          `json:"number" db:"account_id"`
	Status     string          `json:"status" db:"status"`
	Accrual    decimal.Decimal `json:"accrual,omitempty" db:"accrual"`
	UploadedAt time.Time       `json:"uploaded_at" db:"uploaded_at"`
}

type ResponseOrder struct {
	Number     string `json:"number" db:"account_id"`
	Status     string `json:"status" db:"status"`
	Accrual    int64  `json:"accrual,omitempty" db:"accrual"`
	UploadedAt string `json:"uploaded_at" db:"uploaded_at"`
}
type UserOrders struct {
	Orders []Order
}

const (
	OrderAccrualStatusINVALID   = "INVALID"
	OrderAccrualStatusPROCESSED = "PROCESSED"
)

type AccrualServiceResponse struct {
	OrderID string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}

type UserBalanceResponse struct {
	Current   float32 `json:"current"`
	Withdrawn float32 `json:"withdrawn"`
}

type WithdrawnRequest struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}

type Withdrawal struct {
	Order       string  `json:"order"`
	Sum         float32 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
