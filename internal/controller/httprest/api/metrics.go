package api

import (
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
	for key, values := range metrics.Counters {
		fmt.Fprintf(w, "- %s: %d<br>", key, values[0])
	}

	fmt.Fprint(w, "</body></html>")
}

func (api *MetricsAPI) GetMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	mname := chi.URLParam(r, "metric")

	if mtype == metrics.CounterMetricType {
		value, err := api.metrics.GetCounterSum(mname)
		if err != nil {
			switch err {
			case metrics.ErrUnknownMetric:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d", value)
		return
	}

	if mtype == metrics.GaugeMetricType {
		value, err := api.metrics.GetGauge(mname)
		if err != nil {
			switch err {
			case metrics.ErrUnknownMetric:
				http.Error(w, err.Error(), http.StatusNotFound)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", strconv.FormatFloat(value, 'f', -1, 64))
		return
	}

	http.Error(w, metrics.ErrUnknownMetricType.Error(), http.StatusBadRequest)
}

func (api *MetricsAPI) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	mname := chi.URLParam(r, "metric")
	mvalue := chi.URLParam(r, "value")

	if mtype == metrics.CounterMetricType {
		counter, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := api.metrics.AppendCounter(mname, counter); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if mtype == metrics.GaugeMetricType {
		gauge, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := api.metrics.UpdateGauge(mname, gauge); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Error(w, metrics.ErrUnknownMetricType.Error(), http.StatusBadRequest)
}
