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
	r.Route("/api/user/", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/json"))
		r.Post("/register", h.RegistrationUser)
		r.Post("/login", h.AuthUser)
	})
	r.Route("/api/", func(r chi.Router) {
		r.Use(Auth)

		r.Get("/user/balance", h.RegistrationUser)
		r.Post("/user/balance/withdraw", h.RegistrationUser)
		r.Get("/user/withdrawals", h.RegistrationUser)
	})
	r.Route("/api/user/orders", func(r chi.Router) {
		r.Use(Auth)
		r.Use(middleware.AllowContentType("text/plain"))
		r.Post("/", h.LoadOrder)
		r.Get("/", h.RegistrationUser)

		r.Get("/test", h.Test)
	})
	return r
}
