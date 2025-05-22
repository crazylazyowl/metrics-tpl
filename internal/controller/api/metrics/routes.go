package metrics

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(metrics *metrics.Usecase) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	api := NewAPI(metrics)
	r.Post("/update/{metric_type}/{metric_name}/{metric_value}", api.Update)
	return r
}
