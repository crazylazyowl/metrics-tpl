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

	r.Group(func(r chi.Router) {
		r.Use(middleware.Compress)
		r.Get("/", api.GetMetrics)
	})

	// text/plain routes
	r.Group(func(r chi.Router) {
		r.Get("/value/{type}/{metric}", api.GetMetric)
		r.Post("/update/{type}/{metric}/{value}", api.UpdateMetric)
	})

	// json routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.JSONContentType)
		r.Use(middleware.Compress)
		r.Use(middleware.Decompress)
		r.Post("/value/", api.GetMetricJSON)
		r.Post("/update/", api.UpdateMetricJSON)
	})

	return r
}
