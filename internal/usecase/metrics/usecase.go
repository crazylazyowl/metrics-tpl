package metrics

import (
	"context"
	"time"
)

type MetricFetcher interface {
	FetchOne(ctx context.Context, metric Metric) (Metric, error)
	Fetch(ctx context.Context) ([]Metric, error)
}

type MetricUpdater interface {
	UpdateOne(ctx context.Context, metric Metric) error
	Update(ctx context.Context, metrics []Metric) error
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

func (u *MetricUsecase) Metric(ctx context.Context, metric Metric) (Metric, error) {
	if metric.ID == "" {
		return Metric{}, ErrMetricEmptyID
	}
	switch metric.Type {
	case Counter, Gauge:
	default:
		return Metric{}, ErrMetricUnknownType
	}
	return u.reg.FetchOne(ctx, metric)
}

func (u *MetricUsecase) Metrics(ctx context.Context) ([]Metric, error) {
	metrics, err := u.reg.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (u *MetricUsecase) UpdateOne(ctx context.Context, metric Metric) error {
	if err := metric.Validate(); err != nil {
		return err
	}
	var err error
	delay := 1
	for range 3 {
		if err = u.reg.UpdateOne(ctx, metric); err != nil {
			time.Sleep(time.Duration(delay) * time.Second)
			delay += 2
			continue
		}
		break
	}
	return err
}

func (u *MetricUsecase) Update(ctx context.Context, metrics []Metric) error {
	if len(metrics) == 0 {
		return nil
	}
	for _, metric := range metrics {
		if err := metric.Validate(); err != nil {
			return err
		}
	}
	var err error
	delay := 1
	for range 3 {
		if err = u.reg.Update(ctx, metrics); err != nil {
			time.Sleep(time.Duration(delay) * time.Second)
			delay += 2
			continue
		}
		break
	}
	return err
}
