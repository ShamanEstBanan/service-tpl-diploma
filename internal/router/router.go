package router

import (
	"github.com/go-chi/chi"
	"service-tpl-diploma/internal/handler"
)

func New(h *handler.Handler) chi.Router {
	r := chi.NewRouter()
	r.MethodNotAllowed(h.RegistrationUser)
	r.Route("/", func(r chi.Router) {
		r.Post("/registration", h.RegistrationUser)
	})
	return r
}
