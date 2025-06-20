package api

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/middleware"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
)

func NewMetricsRouter(metrics *metrics.Usecase) http.Handler {
	api := NewMetricsAPI(metrics)

	r := chi.NewRouter()

	r.With(middleware.Compress).Get("/", api.GetMetrics)

	r.Get("/value/{type}/{metric}", api.GetMetric)
	r.Post("/update/{type}/{metric}/{value}", api.UpdateMetric)

	r.With(middleware.JSONContentType, middleware.Compress).
		Group(func(r chi.Router) {
			r.Post("/value/", api.GetMetricJSON)
			r.Post("/update/", api.UpdateMetricJSON)
		})

	return r
}
