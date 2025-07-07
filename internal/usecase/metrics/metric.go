package metrics

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
)

type Metric struct {
	ID      string     `json:"id"`
	Type    MetricType `json:"type"`
	Counter *int64     `json:"delta,omitempty"`
	Gauge   *float64   `json:"value,omitempty"`
}

func (metric Metric) Validate() error {
	if metric.ID == "" {
		return ErrEmptyMetricID
	}
	switch metric.Type {
	case Counter:
		if metric.Counter == nil {
			return ErrInvalidCounterValue
		}
	case Gauge:
		if metric.Gauge == nil {
			return ErrInvalidGaugeValue
		}
	default:
		return ErrUnknownMetricType
	}
	return nil
}
