package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
)

type MetricsAPI struct {
	metrics *metrics.Usecase
}

func NewMetricsAPI(metrics *metrics.Usecase) *MetricsAPI {
	return &MetricsAPI{metrics: metrics}
}

func (api *MetricsAPI) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := api.metrics.GetMetrics()

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "<html><head></head><body>")

	fmt.Fprint(w, "Gauge: <br>")
	keys := make([]string, 0, len(metrics.Gauges))
	for key := range metrics.Gauges {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(w, "- %s: %s<br>", key, strconv.FormatFloat(metrics.Gauges[key], 'f', -1, 64))
	}

	fmt.Fprint(w, "Counter: <br>")
	for key, value := range metrics.Counters {
		fmt.Fprintf(w, "- %s: %d<br>", key, value)
	}

	fmt.Fprint(w, "</body></html>")
}

func (api *MetricsAPI) GetMetric(w http.ResponseWriter, r *http.Request) {
	metric := metrics.Metric{
		ID:   chi.URLParam(r, "metric"),
		Type: chi.URLParam(r, "type"),
	}
	metric, err := api.metrics.GetMetric(metric)
	if err != nil {
		switch {
		case errors.As(err, &metrics.ErrMetricUsecase{}): // TODO: handle this error properly
			errNotFound(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	switch metric.Type {
	case metrics.CounterMetricType:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d", *metric.Counter)
	case metrics.GaugeMetricType:
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", strconv.FormatFloat(*metric.Gauge, 'f', -1, 64))
	}
}

func (api *MetricsAPI) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	metric := metrics.Metric{
		ID:   chi.URLParam(r, "metric"),
		Type: chi.URLParam(r, "type"),
	}
	value := chi.URLParam(r, "value")
	switch metric.Type {
	case metrics.CounterMetricType:
		counter, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			errBadRequest(w, err)
			return
		}
		metric.Counter = &counter
	case metrics.GaugeMetricType:
		gauge, err := strconv.ParseFloat(value, 64)
		if err != nil {
			errBadRequest(w, err)
			return
		}
		metric.Gauge = &gauge
	default:
		errBadRequest(w, metrics.ErrUnknownMetricType)
		return
	}
	if err := api.metrics.UpdateMetric(metric); err != nil {
		switch err {
		case metrics.ErrUnknownMetricID:
			errNotFound(w, err)
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
	metric, err := api.metrics.GetMetric(metric)
	if err != nil {
		switch {
		case errors.As(err, &metrics.ErrMetricUsecase{}): // TODO: handle this error properly
			errNotFound(w, err)
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
	if err := api.metrics.UpdateMetric(metric); err != nil {
		switch {
		case errors.As(err, &metrics.ErrMetricUsecase{}): // TODO: handle this error properly
			errBadRequest(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	metric, err := api.metrics.GetMetric(metric)
	if err != nil {
		switch {
		case errors.As(err, &metrics.ErrMetricUsecase{}): // TODO: handle this error properly
			errNotFound(w, err)
		default:
			errInternalServerError(w, err)
		}
		return
	}
	writeJSON(w, http.StatusOK, &metric)
}
