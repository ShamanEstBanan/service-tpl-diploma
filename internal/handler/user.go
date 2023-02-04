package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/errs"
)

func (h *Handler) RegistrationUser(w http.ResponseWriter, r *http.Request) {
	newUser := domain.NewUser{}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		h.lg.Error("Error Incorrect json: ", zap.Any("err", err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	authToken, err := h.service.CreateUser(r.Context(), newUser)
	if errors.Is(err, errs.ErrLoginIsEmpty) {
		http.Error(w, errs.ErrLoginIsEmpty.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, errs.ErrPasswordIsEmpty) {
		http.Error(w, errs.ErrPasswordIsEmpty.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, errs.ErrLoginAlreadyExist) {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Login already exist", http.StatusConflict)
		return
	}

	if err != nil {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", authToken)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Success create"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}

func (h *Handler) AuthUser(w http.ResponseWriter, r *http.Request) {
	user := domain.AuthUser{}
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Incorrect json", http.StatusBadRequest)
		return
	}
	token, err := h.service.AuthUser(r.Context(), user)
	if errors.Is(err, errs.ErrLoginIsEmpty) {
		http.Error(w, errs.ErrLoginIsEmpty.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, errs.ErrPasswordIsEmpty) {
		http.Error(w, errs.ErrPasswordIsEmpty.Error(), http.StatusBadRequest)
		return
	}

	if errors.Is(err, errs.ErrInvalidLoginOrPassword) {
		http.Error(w, errs.ErrInvalidLoginOrPassword.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Authorization", token)
	_, err = w.Write([]byte("Success auth"))
	if err != nil {
		h.lg.Error(err.Error())
		return
	}
}
