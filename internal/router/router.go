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
		//TODO рефакторинг чреез r.With
		//r.With(middleware.Logger).Post()
	})
	r.Route("/api/", func(r chi.Router) {
		r.Use(Auth)

		r.Post("/user/balance/withdraw", h.MakeWithdraw)
		r.Get("/user/withdrawals", h.GetHistoryWithdrawals)
	})
	r.Route("/api/user/orders", func(r chi.Router) {
		r.Use(Auth)
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/", h.LoadOrder)
		r.Get("/", h.GetUserOrders)

		r.Get("/test", h.Test)
	})
	return r
}
