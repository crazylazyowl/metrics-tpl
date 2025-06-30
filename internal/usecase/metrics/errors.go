package metrics

import (
	"fmt"
)

type ErrUnknownMetric struct{}

func (e ErrUnknownMetric) Error() string {
	return "metric error"
}

var (
	ErrUnknownMetricID = fmt.Errorf("%w: unknown metric id", ErrUnknownMetric{})
)

type ErrInvalidMetric struct{}

func (e ErrInvalidMetric) Error() string {
	return "invalid metric"
}

var (
	ErrUnknownMetricType   = fmt.Errorf("%w: unknown metric type", ErrInvalidMetric{})
	ErrEmptyMetricID       = fmt.Errorf("%w: metric id is empty", ErrInvalidMetric{})
	ErrInvalidCounterValue = fmt.Errorf("%w: metric counter value is invalid", ErrInvalidMetric{})
	ErrInvalidGaugeValue   = fmt.Errorf("%w: metric gauge value is invalid", ErrInvalidMetric{})
)
