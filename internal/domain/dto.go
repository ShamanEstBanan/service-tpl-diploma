package domain

import (
	"context"
	"github.com/shopspring/decimal"
	"time"
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
	OrderInternalStatusINVALID    = "INVALID"
	OrderInternalStatusPROCESSED  = "PROCESSED"
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
	OrderAccrualStatusREGISTERED = "REGISTERED"
	OrderAccrualStatusINVALID    = "INVALID"
	OrderAccrualStatusPROCESSING = "PROCESSING"
	OrderAccrualStatusPROCESSED  = "PROCESSED"
)

type AccrualServiceResponse struct {
	OrderId string          `json:"order"`
	Status  string          `json:"status"`
	Accrual decimal.Decimal `json:"accrual,omitempty"`
}
