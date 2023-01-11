package handler

import "net/http"

func (h *Handler) RegistrationUser(w http.ResponseWriter, r *http.Request) {

	_, err := w.Write([]byte("Success"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
