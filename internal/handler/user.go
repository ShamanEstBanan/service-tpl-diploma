package handler

import (
	"encoding/json"
	"errors"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"net/http"
	"service-tpl-diploma/internal/app/domain"
	"service-tpl-diploma/internal/errs"
)

var (
	ErrInvalidLoginOrPassword = errors.New("login or password is invalid")
	ErrLoginIsEmpty           = errors.New("login is empty")
	ErrPasswordIsEmpty        = errors.New("password is empty")
)

func (h *Handler) RegistrationUser(w http.ResponseWriter, r *http.Request) {
	newUser := domain.NewUser{}
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Incorrect json", http.StatusInternalServerError)
		return
	}

	if newUser.Login == "" {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, ErrLoginIsEmpty.Error(), http.StatusBadRequest)
		return
	} else if newUser.Password == "" {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, ErrPasswordIsEmpty.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.CreateUser(r.Context(), newUser)

	duplicateErr := errs.NewSQLError(pgerrcode.DuplicateJSONObjectKeyValue)
	if errors.As(err, &duplicateErr) {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Login already exist", http.StatusConflict)
		return
	}

	if err != nil {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

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
	if errors.As(err, &ErrLoginIsEmpty) {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, "Login is empty", http.StatusBadRequest)
		return
	}

	if errors.As(err, &ErrPasswordIsEmpty) {
		h.lg.Error("Error: ", zap.Any("err", err))
		http.Error(w, ErrPasswordIsEmpty.Error(), http.StatusBadRequest)
		return
	}

	if errors.As(err, &ErrInvalidLoginOrPassword) {
		http.Error(w, ErrInvalidLoginOrPassword.Error(), http.StatusUnauthorized)
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
