package httprest

import (
	"net/http"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/api"
	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/middleware"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
)

func NewRouter(metrics *metrics.Usecase) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", api.NewMetricsRouter(metrics))
	return r
}
