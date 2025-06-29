package metrics

import (
	"fmt"
)

type ErrMetricUsecase struct{}

func (e ErrMetricUsecase) Error() string {
	return "metric error"
}

var (
	ErrUnknownMetricID   = fmt.Errorf("%w: unknown metric id", ErrMetricUsecase{})
	ErrUnknownMetricType = fmt.Errorf("%w: unknown metric type", ErrMetricUsecase{})
	ErrMetricEmptyID     = fmt.Errorf("%w: metric id is empty", ErrMetricUsecase{})
	ErrMetricValue       = fmt.Errorf("%w: metric value is invalid", ErrMetricUsecase{})
)
