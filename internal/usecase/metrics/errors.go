package metrics

import "errors"

var (
	ErrUnknownMetric     = errors.New("unknown metric")
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrBadCounterValue   = errors.New("bad counter value, integer is expected")
	ErrBadGaugeValue     = errors.New("bad gauge value, float is expected")
)
