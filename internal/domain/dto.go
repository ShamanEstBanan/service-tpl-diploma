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
	Number     string    `json:"number"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
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
