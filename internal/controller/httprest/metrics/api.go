package metrics

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"

	"github.com/go-chi/chi/v5"
)

type API struct {
	metrics *metrics.Usecase
}

func NewAPI(metrics *metrics.Usecase) *API {
	return &API{metrics: metrics}
}

func (api *API) GetMetrics(w http.ResponseWriter, r *http.Request) {
	counters, gauges := api.metrics.Metrics()

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, "<html><head></head><body>")

	fmt.Fprint(w, "Gauge: <br>")
	keys := make([]string, 0, len(gauges))
	for key := range gauges {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(w, "- %s: %s<br>", key, gauges[key].String())
	}

	fmt.Fprint(w, "Counter: <br>")
	for key, values := range counters {
		fmt.Fprintf(w, "- %s: %d<br>", key, values[0])
	}

	fmt.Fprint(w, "</body></html>")
}

func (api *API) GetMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	mname := chi.URLParam(r, "metric")

	if mtype == metrics.CounterMetricType {
		value, err := api.metrics.CounterSum(mname)
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
		value, err := api.metrics.Gauge(chi.URLParam(r, "metric"))
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
		fmt.Fprintf(w, value.String())
		return
	}

	http.Error(w, metrics.ErrUnknownMetricType.Error(), http.StatusBadRequest)
}

func (api *API) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	mtype := chi.URLParam(r, "type")
	mname := chi.URLParam(r, "metric")
	mvalue := chi.URLParam(r, "value")

	if mtype == metrics.CounterMetricType {
		counter, err := metrics.CounterFromString(mvalue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := api.metrics.UpdateCounter(mname, counter); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	if mtype == metrics.GaugeMetricType {
		gauge, err := metrics.GaugeFromString(mvalue)
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
