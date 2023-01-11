package handler

import (
	"go.uber.org/zap"
	"service-tpl-diploma/internal/storage"
)

type Handler struct {
	lg *zap.Logger
	st *storage.Storage
}

func New(lg *zap.Logger, st *storage.Storage) *Handler {
	return &Handler{
		lg: lg,
		st: st,
	}
}
