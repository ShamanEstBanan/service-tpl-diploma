package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"service-tpl-diploma/internal/errs"
)

func (h *Handler) LoadOrder(w http.ResponseWriter, r *http.Request) {
	reader := r.Body
	body, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println(err)
		return
	}

	on := string(body)
	orderNumber, err := strconv.Atoi(on)
	if err != nil {
		http.Error(w, "Invalid format of order number", http.StatusUnprocessableEntity)
		return
	}
	h.lg.Sugar().Info("Order number: ", orderNumber)
	userID := r.Header.Get("userID")
	err = h.service.LoadOrder(r.Context(), orderNumber, userID)
	if errors.Is(err, errs.ErrOrderAlreadyExist) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	if errors.Is(err, errs.ErrOrderAlreadyUploaded) {
		http.Error(w, err.Error(), http.StatusOK)
		return
	}
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte("Order successful loaded to system"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}

func (h *Handler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID")

	orders, err := h.service.GetUserOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		_, err = w.Write([]byte(""))
		if err != nil {
			h.lg.Error(err.Error())
			return
		}
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(orders)
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
