package metrics

import "errors"

var (
	ErrUnknownMetric     = errors.New("unknown metric")
	ErrUnknownMetricType = errors.New("unknown metric type")
	ErrBadMetricValue    = errors.New("bad metric value")
)
