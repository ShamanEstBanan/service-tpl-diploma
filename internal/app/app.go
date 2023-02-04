package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"service-tpl-diploma/internal/app/migrate"
	"service-tpl-diploma/internal/config"
	"service-tpl-diploma/internal/domain"
	"service-tpl-diploma/internal/handler"
	"service-tpl-diploma/internal/router"
	"service-tpl-diploma/internal/storage"
	"service-tpl-diploma/internal/usecase"
	"service-tpl-diploma/internal/workers"
	"service-tpl-diploma/pkg/logger"
)

const (
	migrationsPath = "internal/app/migrate/migrations/"
)

type App struct {
	logger     *zap.Logger
	HTTPServer *http.Server
}

// listen OS signals to gracefully shutdown server
func listen(ctx context.Context, server *http.Server) error {
	srv := server

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(srv.ListenAndServe)
	log.Println("server started 1", zap.String("addr:", server.Addr))
	eg.Go(func() error {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			fmt.Println("Shutdown error:", err)
		}
		return err
	})

	return eg.Wait()
}

// New return app instance
func New() (*App, error) {
	cfg := config.New()

	lg, err := logger.New(cfg.Debug)
	if err != nil {
		return nil, err
	}

	// Init migrations
	if err := migrate.Run(cfg.PostgresDSN, migrate.WithPath(migrationsPath)); err != nil {
		log.Fatalf("failed executing migrate DB: %v\n", err) //nolint: revive
	}
	ctx := context.Background()

	// init storage
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}
	st := storage.New(pool, lg)

	// workerpool
	jobsCh := make(chan domain.Job, cfg.JobCount)
	go func() {
		ctx, cancel := signal.NotifyContext(
			context.Background(),
			syscall.SIGINT,
			syscall.SIGTERM,
		)
		defer cancel()
		err = workers.RunPool(ctx, cfg.Concurrency, jobsCh)
		if err != nil {
			lg.Error(" ERROR worker-pool:", zap.Error(err))
		}
	}()

	// init worker updater
	w := workers.NewStatusUpdater(st, jobsCh, lg, cfg.AccrualSystemAddress)
	w.Start()

	// init service
	usecases := usecase.New(lg, st)
	// init server
	h := handler.New(lg, usecases)
	r := router.New(h)
	server := &http.Server{
		Addr:              cfg.RunAddress,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}
	return &App{
		logger:     lg,
		HTTPServer: server,
	}, nil
}

// Run starts our server
func (app *App) Run() error {
	// gracefully shutdown
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return listen(ctx, app.HTTPServer)
	})
	err := eg.Wait()
	if err != nil {
		return err
	}
	return nil
}
