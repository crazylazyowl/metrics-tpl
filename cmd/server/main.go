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
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	conf, err := loadConfig()
	if err != nil {
		log.Err(err).Msg("faild to load config")
		return
	}

	ctx := context.Background()

	notifyCtx, notifyCancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer notifyCancel()

	storage := memstorage.New()
	if conf.restore {
		if err := storage.Restore(conf.fileStoragePath); err != nil {
			log.Err(err).Msg("failed to restore memstorage")
			// return
		}
	}
	go storage.Backup(notifyCtx, conf.fileStoragePath, conf.storeInterval)

	usecase := metrics.New(storage)

	router := httprest.NewRouter(usecase)

	server := http.Server{
		Addr:    conf.address,
		Handler: router,
	}

	go func() {
		<-notifyCtx.Done()
		log.Debug().Msg("shutdown")
		timeoutCtx, timeoutCancel := context.WithTimeout(ctx, 10*time.Second)
		defer timeoutCancel()
		server.Shutdown(timeoutCtx)
	}()

	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Err(err).Msg("listen and server failed")
		}
	}
}
