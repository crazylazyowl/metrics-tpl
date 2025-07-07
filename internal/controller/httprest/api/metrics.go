package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/crazylazyowl/metrics-tpl/internal/controller/httprest/middleware"
	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
)

type MetricsAPI struct {
	metrics *metrics.MetricUsecase
}

func NewMetricsAPI(metrics *metrics.MetricUsecase) *MetricsAPI {
	return &MetricsAPI{metrics: metrics}
}

func NewMetricsRouter(metrics *metrics.MetricUsecase) http.Handler {
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
		r.Post("/updates/", api.UpdateMetricsJSON)
	})

	return r
}

func (api *MetricsAPI) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metricList, err := api.metrics.Metrics(r.Context())
	if err != nil {
		errInternalServerError(w, err) // TODO: review error handling
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "<html><head></head><body>")

	for _, m := range metricList {
		switch m.Type {
		case metrics.Counter:
			fmt.Fprintf(w, "- %s: %d<br>", m.ID, *m.Counter)
		case metrics.Gauge:
			fmt.Fprintf(w, "- %s: %s<br>", m.ID, strconv.FormatFloat(*m.Gauge, 'f', -1, 64))
		}
	}

	fmt.Fprint(w, "</body></html>")
}

func (api *MetricsAPI) GetMetric(w http.ResponseWriter, r *http.Request) {
	metric := metrics.Metric{
		ID:   chi.URLParam(r, "metric"),
		Type: metrics.MetricType(chi.URLParam(r, "type")),
	}
	metric, err := api.metrics.Metric(r.Context(), metric)
	if err != nil {
		switch {
		case errors.Is(err, metrics.ErrMetricNotFound):
			errNotFound(w, err)
		case errors.Is(err, metrics.ErrMetricInvalid):
			errBadRequest(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	switch metric.Type {
	case metrics.Counter:
		fmt.Fprintf(w, "%d", *metric.Counter)
	case metrics.Gauge:
		fmt.Fprintf(w, "%s", strconv.FormatFloat(*metric.Gauge, 'f', -1, 64))
	}
}

func (api *MetricsAPI) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := metrics.Metric{
		ID:   chi.URLParam(r, "metric"),
		Type: metrics.MetricType(chi.URLParam(r, "type")),
	}
	value := chi.URLParam(r, "value")
	switch metric.Type {
	case metrics.Counter:
		counter, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			errBadRequest(w, err)
			return
		}
		metric.Counter = &counter
	case metrics.Gauge:
		gauge, err := strconv.ParseFloat(value, 64)
		if err != nil {
			errBadRequest(w, err)
			return
		}
		metric.Gauge = &gauge
	default:
		errBadRequest(w, metrics.ErrMetricUnknownType)
		return
	}
	if err := api.metrics.UpdateOne(r.Context(), metric); err != nil {
		switch {
		case errors.Is(err, metrics.ErrMetricInvalid):
			errBadRequest(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *MetricsAPI) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metric
	if err := readJSON(r, &metric); err != nil {
		errBadRequest(w, err)
		return
	}
	metric, err := api.metrics.Metric(r.Context(), metric)
	if err != nil {
		switch {
		case errors.Is(err, metrics.ErrMetricNotFound):
			errNotFound(w, err)
		case errors.Is(err, metrics.ErrMetricInvalid):
			errBadRequest(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, &metric)
}

func (api *MetricsAPI) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metric metrics.Metric
	if err := readJSON(r, &metric); err != nil {
		errBadRequest(w, err)
		return
	}
	ctx := r.Context()
	// TODO: race condition between UpdateMetric and GetMetric
	if err := api.metrics.UpdateOne(ctx, metric); err != nil {
		switch {
		case errors.Is(err, metrics.ErrMetricInvalid):
			errBadRequest(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	metric, err := api.metrics.Metric(ctx, metric)
	if err != nil {
		switch {
		case errors.Is(err, metrics.ErrMetricInvalid):
			errNotFound(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, &metric)
}

func (api *MetricsAPI) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) {
	var many []metrics.Metric
	if err := readJSON(r, &many); err != nil {
		errBadRequest(w, err)
	}
	ctx := r.Context()
	if err := api.metrics.Update(ctx, many); err != nil {
		switch {
		case errors.Is(err, metrics.ErrMetricInvalid):
			errBadRequest(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}
