package metrics

import (
	"net/http"
	"strings"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

type API struct {
	metrics *metrics.Usecase
}

func NewAPI(metrics *metrics.Usecase) *API {
	return &API{metrics: metrics}
}

func (api *API) Update(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	if err := api.metrics.Update(parts[0], parts[1], parts[2]); err != nil {
		switch err {
		case metrics.ErrUnknownMetric:
			http.Error(w, err.Error(), http.StatusNotFound)
		case metrics.ErrUnknownMetricType, metrics.ErrBadMetricValue:
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(parts, "|")))
}
