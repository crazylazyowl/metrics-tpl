package metrics

import "context"

type MetricFetcher interface {
	Fetch(ctx context.Context) ([]Metric, error)
	FetchOne(ctx context.Context, m Metric) (Metric, error)
}

type MetricUpdater interface {
	Update(ctx context.Context, m Metric) error
}

type MetricRegistry interface {
	MetricFetcher
	MetricUpdater
}

type MetricUsecase struct {
	reg MetricRegistry
}

func New(reg MetricRegistry) *MetricUsecase {
	return &MetricUsecase{reg: reg}
}

func (u *MetricUsecase) Metrics(ctx context.Context) ([]Metric, error) {
	metrics, err := u.reg.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (u *MetricUsecase) Metric(ctx context.Context, m Metric) (Metric, error) {
	if m.ID == "" {
		return Metric{}, ErrEmptyMetricID
	}
	switch m.Type {
	case CounterMetricType, GaugeMetricType:
	default:
		return Metric{}, ErrUnknownMetricType
	}
	metric, err := u.reg.FetchOne(ctx, m)
	if err != nil {
		return Metric{}, err
	}
	return metric, nil
}

func (u *MetricUsecase) Update(ctx context.Context, m Metric) error {
	if m.ID == "" {
		return ErrEmptyMetricID
	}
	switch m.Type {
	case CounterMetricType:
		if m.Counter == nil {
			return ErrInvalidCounterValue
		}
	case GaugeMetricType:
		if m.Gauge == nil {
			return ErrInvalidGaugeValue
		}
	default:
		return ErrUnknownMetricType
	}
	if err := u.reg.Update(ctx, m); err != nil {
		return err
	}
	return nil
}
