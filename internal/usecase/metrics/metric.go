package metrics

const (
	CounterMetricType = "counter"
	GaugeMetricType   = "gauge"
)

type Metric struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Counter *int64   `json:"delta,omitempty"`
	Gauge   *float64 `json:"value,omitempty"`
}
