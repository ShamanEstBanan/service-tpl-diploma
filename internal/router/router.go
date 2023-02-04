package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"service-tpl-diploma/internal/handler"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(1))
	r.Route("/api/user/", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/register", h.RegistrationUser)
		r.Post("/login", h.AuthUser)
		r.With(Auth).Get("/balance", h.GetBalance)
		r.With(Auth).Post("/balance/withdraw", h.MakeWithdraw)
		r.With(Auth).Get("/withdrawals", h.GetHistoryWithdrawals)
		// TODO рефакторинг чреез r.With
	})
	r.Route("/api/user/orders", func(r chi.Router) {
		r.Use(Auth)
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/", h.LoadOrder)
		r.Get("/", h.GetUserOrders)
	})
	r.Route("/admin", func(r chi.Router) {
		r.Get("/test", h.Test)
	})
	return r
}
