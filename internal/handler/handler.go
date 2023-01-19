package handler

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"service-tpl-diploma/internal/app/domain"
)

type service interface {
	CreateUser(ctx context.Context, user domain.NewUser) error
	AuthUser(ctx context.Context, User domain.AuthUser) (string, error)
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
