package metrics

import (
	"errors"
	"fmt"
)

type errMetricInvalid struct{}

func (e *errMetricInvalid) Error() string {
	return "metric is invalid"
}

var (
	ErrMetricNotFound       = errors.New("metric not found")
	ErrMetricInvalid        = &errMetricInvalid{}
	ErrMetricUnknownType    = fmt.Errorf("%w: unknown type", ErrMetricInvalid)
	ErrMetricEmptyID        = fmt.Errorf("%w: id is empty", ErrMetricInvalid)
	ErrMetricInvalidCounter = fmt.Errorf("%w: counter value is invalid", ErrMetricInvalid)
	ErrMetricInvalidGauge   = fmt.Errorf("%w: gauge value is invalid", ErrMetricInvalid)
)
