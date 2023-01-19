package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"service-tpl-diploma/internal/handler"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Route("/api/user/", func(r chi.Router) {
		r.Post("/register", h.RegistrationUser) //Content-Type: application/json
		r.Post("/login", h.AuthUser)            // Content-Type: application/json
	})
	r.Route("/api/user/", func(r chi.Router) {
		r.Use(Auth)
		r.Post("/orders", h.RegistrationUser)           //Content-Type: text/plain
		r.Get("/orders", h.RegistrationUser)            //-
		r.Get("/balance", h.RegistrationUser)           // -
		r.Post("/balance/withdraw", h.RegistrationUser) //Content-Type: application/json
		r.Get("/withdrawals", h.RegistrationUser)       // -
	})
	r.Route("/test", func(r chi.Router) {
		r.Use(Auth)
		r.Get("/test", h.Test)
	})
	return r
}
