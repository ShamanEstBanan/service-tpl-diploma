package handler

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")
	balanceInfo, err := h.service.GetUserBalance(r.Context(), userID)
	if errors.Is(err, errs.ErrNoWithdrawn) {
		w.Header().Set("content-type", "application/json")
		_ = json.NewEncoder(w).Encode(balanceInfo)
		return
	}
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
	userID := r.Header.Get("userID")
	withdrawn := domain.WithdrawnRequest{}
	err := json.NewDecoder(r.Body).Decode(&withdrawn)
	if err != nil {
		h.lg.Error("Error: ", zap.Error(err))
		http.Error(w, "Incorrect json", http.StatusBadRequest)
		return
	}
	err = h.service.MakeWithdrawn(r.Context(), userID, withdrawn.Order, withdrawn.Sum)
	if errors.Is(err, errs.ErrWithdrawnAlreadyDoneForThisOrder) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if errors.Is(err, errs.ErrNotEnoughtPoints) {
		http.Error(w, err.Error(), http.StatusPaymentRequired)
		return
	}
	if errors.Is(err, errs.ErrInvalidOrderID) {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	if err != nil {
		h.lg.Error("Error: ", zap.Error(err))
		http.Error(w, "Incorrect json", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
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
