package metrics

type Metric struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Counter *int64   `json:"delta,omitempty"`
	Gauge   *float64 `json:"value,omitempty"`
}

func (m Metric) Validate() error {
	if m.ID == "" {
		return ErrMetricEmptyID
	}
	switch m.Type {
	case CounterMetricType, GaugeMetricType:
	default:
		return ErrUnknownMetricType
	}
	return nil
}
