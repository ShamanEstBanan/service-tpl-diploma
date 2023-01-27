package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"service-tpl-diploma/internal/domain"
)

type service interface {
	CreateUser(ctx context.Context, user domain.NewUser) error
	AuthUser(ctx context.Context, User domain.AuthUser) (string, error)
	LoadOrder(ctx context.Context, orderId int, userId string) error
	GetUserOrders(ctx context.Context, userID string) (orders []domain.ResponseOrder, err error)
	GetUserBalance(ctx context.Context, userID string) (balance domain.UserBalanceResponse, err error)
	MakeWithdrawn(ctx context.Context, userID string, orderID string, sum float32) (err error)
}
type Handler struct {
	lg      *zap.Logger
	service service
}

func New(lg *zap.Logger, service service) *Handler {
	return &Handler{
		lg:      lg,
		service: service,
	}
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Header.Get("userId")
	fmt.Println("userLogin=", userLogin)
	_, err := w.Write([]byte("Test succses"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
