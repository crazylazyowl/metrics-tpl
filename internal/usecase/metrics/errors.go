package metrics

import (
	"errors"
	"fmt"
)

type ErrInvalidMetric struct{}

func NewInvalidMetricError() *ErrInvalidMetric {
	return &ErrInvalidMetric{}
}

func (e *ErrInvalidMetric) Error() string {
	return "invalid metric"
}

var (
	ErrNotFound            = errors.New("metric not found")
	ErrUnknownMetricType   = fmt.Errorf("%w: unknown type", NewInvalidMetricError())
	ErrEmptyMetricID       = fmt.Errorf("%w: id is empty", NewInvalidMetricError())
	ErrInvalidCounterValue = fmt.Errorf("%w: counter value is invalid", NewInvalidMetricError())
	ErrInvalidGaugeValue   = fmt.Errorf("%w: gauge value is invalid", NewInvalidMetricError())
)
