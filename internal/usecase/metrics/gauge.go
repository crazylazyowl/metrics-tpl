package metrics

import "strconv"

const GaugeMetricType = "gauge"

type Gauge float64

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func GaugeFromString(value string) (Gauge, error) {
	gauge, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, ErrBadCounterValue
	}
	return Gauge(gauge), nil
}
