package api

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
)

func NewMetricsRouter(metrics *metrics.Usecase) http.Handler {
	api := NewMetricsAPI(metrics)
	r := chi.NewRouter()
	r.Get("/", api.GetMetrics)
	r.Get("/value/{type}/{metric}", api.GetMetric)
	r.Post("/update/{type}/{metric}/{value}", api.UpdateMetric)
	return r
}
