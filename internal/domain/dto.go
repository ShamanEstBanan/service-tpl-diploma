package domain

import "context"

type NewUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

const (
	OrderStatusNEW        = "NEW"
	OrderStatusPROCESSING = "PROCESSING"
	OrderStatusINVALID    = "INVALID"
	OrderStatusPROCESSED  = "PROCESSED"
)

type Job interface {
	Run(ctx context.Context) error
}
