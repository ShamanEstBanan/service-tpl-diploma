package handler

import "net/http"

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	h.lg.Sugar().Info("INFO:", r.Host+r.URL.Path)
	_, err := w.Write([]byte("Success"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}

func (h *Handler) MakeWithdraw(w http.ResponseWriter, r *http.Request) {
	h.lg.Sugar().Info("INFO:", r.Host+r.URL.Path)
	_, err := w.Write([]byte("Success"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}

func (h *Handler) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	h.lg.Sugar().Info("INFO:", r.Host+r.URL.Path)
	_, err := w.Write([]byte("Success"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
