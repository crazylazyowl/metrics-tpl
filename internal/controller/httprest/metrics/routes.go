package metrics

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(metrics *metrics.Usecase) http.Handler {
	api := NewAPI(metrics)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", api.GetMetrics)

	r.Get("/value/{type}/{metric}", api.GetMetric)
	r.Post("/update/{type}/{metric}/{value}", api.UpdateMetric)

	return r
}
