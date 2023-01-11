package storage

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	db *pgxpool.Pool
	lg *zap.Logger
}

func New(db *pgxpool.Pool, lg *zap.Logger) *Storage {
	return &Storage{
		db: db,
		lg: lg,
	}
}
