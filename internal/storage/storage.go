package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type storage struct {
	db *pgxpool.Pool
	lg *zap.Logger
}

func New(db *pgxpool.Pool, lg *zap.Logger) *storage {
	return &storage{
		db: db,
		lg: lg,
	}
}
