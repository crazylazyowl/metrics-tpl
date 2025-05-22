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

	r.Get("/", api.Metrics)

	r.Route("/value", func(r chi.Router) {
		r.Get("/counter/{metric}", api.Counter)
		r.Get("/gauge/{metric}", api.Gauge)
	})

	r.Route("/update", func(r chi.Router) {
		r.Post("/counter/{metric}/{value}", api.UpdateCounter)
		r.Post("/gauge/{metric}/{value}", api.UpdateGauge)
	})

	return r
}
