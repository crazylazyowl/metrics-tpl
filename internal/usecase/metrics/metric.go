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

func (metric Metric) Validate() error {
	if metric.ID == "" {
		return ErrEmptyMetricID
	}
	switch metric.Type {
	case CounterMetricType:
		if metric.Counter == nil {
			return ErrInvalidCounterValue
		}
	case GaugeMetricType:
		if metric.Gauge == nil {
			return ErrInvalidGaugeValue
		}
	default:
		return ErrUnknownMetricType
	}
	return nil
}
