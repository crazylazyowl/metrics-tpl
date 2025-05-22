package metrics

import "errors"

var (
	ErrUnknownMetric   = errors.New("unknown metric")
	ErrBadCounterValue = errors.New("bad counter value, integer is expected")
	ErrBadGaugeValue   = errors.New("bad gauge value, floag is expected")
)
