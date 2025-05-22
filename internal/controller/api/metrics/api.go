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

func (api *API) Metrics(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, "- %s: %.2f<br>", key, gauges[key])
	}

	fmt.Fprint(w, "Counter: <br>")
	for key, values := range counters {
		fmt.Fprintf(w, "- %s: %d<br>", key, values[0])
	}

	fmt.Fprint(w, "</body></html>")
}

func (api *API) Counter(w http.ResponseWriter, r *http.Request) {
	value, err := api.metrics.CounterSum(chi.URLParam(r, "metric"))
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
}

func (api *API) Gauge(w http.ResponseWriter, r *http.Request) {
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
	fmt.Fprintf(w, "%.3f", value)
}

func (api *API) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	counter, err := metrics.CounterFromString(chi.URLParam(r, "value"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := api.metrics.UpdateCounter(chi.URLParam(r, "metric"), counter); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (api *API) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	gauge, err := metrics.GaugeFromString(chi.URLParam(r, "value"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := api.metrics.UpdateGauge(chi.URLParam(r, "metric"), gauge); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
