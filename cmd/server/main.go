package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest"
	"github.com/crazylazyowl/metrics-tpl/internal/repository/memstorage"
	"github.com/crazylazyowl/metrics-tpl/internal/repository/postgres"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logger := log.With().Logger() // TODO: move to separate package

	conf, err := loadConfig()
	if err != nil {
		logger.Err(err).Msg("faild to load config")
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var metricsStorage metrics.MetricRegistry
	var pingStorage ping.Pinger

	if conf.db.dsn == "" {
		logger.Debug().Msg("init memstorage")
		stor, err := memstorage.New(ctx, memstorage.Options{
			Restore:        conf.storage.restore,
			BackupPath:     conf.storage.backupPath,
			BackupInterval: time.Duration(conf.storage.backupInterval) * time.Second,
		})
		if err != nil {
			logger.Err(err).Msg("failed to create memstorage")
			return
		}
		defer stor.Close(ctx)
		metricsStorage = stor
		pingStorage = stor
	} else {
		logger.Debug().Msg("init postgres")
		stor, err := postgres.NewPostgresStorage(ctx, postgres.Options{DSN: conf.db.dsn, Migrations: "file://migrations"})
		if err != nil {
			logger.Err(err).Msg("failed to create postgres storage")
			return
		}
		defer stor.Close(ctx)
		metricsStorage = stor
		pingStorage = stor
	}

	metricsUsecase := metrics.New(metricsStorage)
	pingUsecase := ping.New(pingStorage)
	router := httprest.NewRouter(metricsUsecase, pingUsecase)
	server := http.Server{
		Addr:    conf.address,
		Handler: router,
	}

	go func() {
		<-ctx.Done()

		logger.Debug().Msg("shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			logger.Err(err).Msg("failed to shutdown server")
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Err(err).Msg("listen and server failed")
	}
}
