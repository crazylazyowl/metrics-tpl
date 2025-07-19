package httprest

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/api"
	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/middleware"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/ping"

	"github.com/go-chi/chi/v5"
)

func NewRouter(metrics *metrics.MetricUsecase, ping *ping.PingUsecase, hmacSecret string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.CheckHMAC(hmacSecret))
	r.Mount("/", api.NewMetricsRouter(metrics))
	r.Mount("/ping", api.NewPingRouter(ping))
	return r
}
