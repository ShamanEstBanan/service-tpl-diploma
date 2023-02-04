package migrate

import (
	"database/sql"
	"embed"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

type options struct {
	driverName string
	path       string
}

type OptionsFunc func(opts *options)

func WithPath(path string) OptionsFunc {
	return func(opts *options) {
		opts.path = path
	}
}

func WithDriver(driver string) OptionsFunc {
	return func(opts *options) {
		opts.driverName = driver
	}
}

func WithFs(fs embed.FS) OptionsFunc {
	return func(opts *options) {
		goose.SetBaseFS(fs)
	}
}

func Run(dsn string, opts ...OptionsFunc) error {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}
	curDriverName := "pgx"
	if o.driverName != "" {
		curDriverName = o.driverName
	}
	sqlDB, err := sql.Open(curDriverName, dsn)
	if err != nil {
		return fmt.Errorf("open database connection error: %w ", err)
	}
	defer func() { _ = sqlDB.Close() }()
	curPath := "/Users/yunbaranik/go/src/service-tpl-diploma/internal/app/migrate/migrations/"
	if o.path != "" {
		curPath = o.path
	}
	if err = goose.Up(sqlDB, curPath); err != nil {
		return fmt.Errorf("up migrations: %w", err)
	}
	return nil
}
