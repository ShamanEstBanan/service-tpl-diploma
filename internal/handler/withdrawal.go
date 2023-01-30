package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")
	balanceInfo, err := h.service.GetUserBalance(r.Context(), userID)
	if errors.Is(err, errs.ErrNoWithdrawn) {
		w.Header().Set("content-type", "application/json")
		err = json.NewEncoder(w).Encode(balanceInfo)
		if err != nil {
			h.lg.Error(err.Error())
		}
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
	if errors.Is(err, errs.ErrNotEnoughPoints) {
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
	userID := r.Header.Get("userID")
	withdrawals, err := h.service.GetUserWithdrawals(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		http.Error(w, "No rows", http.StatusNoContent)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
