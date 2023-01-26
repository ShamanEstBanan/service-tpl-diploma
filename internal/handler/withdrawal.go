package handler

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")
	balanceInfo, err := h.service.GetUserBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(balanceInfo)
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

func (h *Handler) GetHistoryWithdrawals(w http.ResponseWriter, r *http.Request) {
	h.lg.Sugar().Info("INFO:", r.Host+r.URL.Path)
	_, err := w.Write([]byte("Success"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
