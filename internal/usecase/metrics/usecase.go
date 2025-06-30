package metrics

import "context"

type MetricsStorage interface {
	GetCounter(ctx context.Context, name string) (int64, error)
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounters(ctx context.Context) map[string]int64
	GetGauges(ctx context.Context) map[string]float64
	UpdateCounter(ctx context.Context, name string, value int64) error
	UpdateGauge(ctx context.Context, name string, value float64) error
}

type Usecase struct {
	storage MetricsStorage
}

func New(repo MetricsStorage) *Usecase {
	return &Usecase{storage: repo}
}

const (
	CounterMetricType = "counter"
	GaugeMetricType   = "gauge"
)

type Metrics struct {
	Counters map[string]int64
	Gauges   map[string]float64
}

func (u *Usecase) GetMetrics(ctx context.Context) Metrics {
	return Metrics{
		Counters: u.storage.GetCounters(ctx),
		Gauges:   u.storage.GetGauges(ctx),
	}
}

// GetMetric returns the metric by its ID and type.
func (u *Usecase) GetMetric(ctx context.Context, m Metric) (Metric, error) {
	if m.ID == "" {
		return Metric{}, ErrEmptyMetricID
	}
	switch m.Type {
	case CounterMetricType:
		value, err := u.storage.GetCounter(ctx, m.ID)
		if err != nil {
			return Metric{}, err
		}
		m.Counter = &value
	case GaugeMetricType:
		value, err := u.storage.GetGauge(ctx, m.ID)
		if err != nil {
			return Metric{}, err
		}
		m.Gauge = &value
	default:
		return Metric{}, ErrUnknownMetricType
	}
	return m, nil
}

// UpdateMetric updates the metric value based on its type and name.
func (u *Usecase) UpdateMetric(ctx context.Context, m Metric) error {
	if m.ID == "" {
		return ErrEmptyMetricID
	}
	switch m.Type {
	case CounterMetricType:
		if m.Counter == nil {
			return ErrInvalidCounterValue
		}
		return u.storage.UpdateCounter(ctx, m.ID, *m.Counter)
	case GaugeMetricType:
		if m.Gauge == nil {
			return ErrInvalidGaugeValue
		}
		return u.storage.UpdateGauge(ctx, m.ID, *m.Gauge)
	}
	return ErrUnknownMetricType
}
