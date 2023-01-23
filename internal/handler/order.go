package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"service-tpl-diploma/internal/errs"
	"strconv"
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

	http.Error(w, err.Error(), http.StatusAccepted)
	_, err = w.Write([]byte("Order successful loaded to system"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}

func (h *Handler) GetOrders(w http.ResponseWriter, r *http.Request) {
	h.lg.Sugar().Info("INFO:", r.Host+r.URL.Path)
	_, err := w.Write([]byte("Success"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
