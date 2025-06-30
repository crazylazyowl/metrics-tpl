package metrics

type MetricsStorage interface {
	GetCounter(name string) (int64, error)
	GetGauge(name string) (float64, error)
	GetCounters() map[string]int64
	GetGauges() map[string]float64
	UpdateCounter(name string, value int64) error
	UpdateGauge(name string, value float64) error
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

func (u *Usecase) GetMetrics() Metrics {
	return Metrics{
		Counters: u.storage.GetCounters(),
		Gauges:   u.storage.GetGauges(),
	}
}

// GetCounterSum returns the value for the specified counter.
func (u *Usecase) GetCounterSum(name string) (int64, error) {
	value, err := u.storage.GetCounter(name)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// GetGauge returnes the value for the specified gauge.
func (u *Usecase) GetGauge(name string) (float64, error) {
	value, err := u.storage.GetGauge(name)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// GetMetric returns the metric by its ID and type.
func (u *Usecase) GetMetric(m Metric) (Metric, error) {
	if m.ID == "" {
		return Metric{}, ErrEmptyMetricID
	}
	switch m.Type {
	case CounterMetricType:
		value, err := u.storage.GetCounter(m.ID)
		if err != nil {
			return Metric{}, err
		}
		m.Counter = &value
	case GaugeMetricType:
		value, err := u.storage.GetGauge(m.ID)
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
func (u *Usecase) UpdateMetric(m Metric) error {
	if m.ID == "" {
		return ErrEmptyMetricID
	}
	switch m.Type {
	case CounterMetricType:
		if m.Counter == nil {
			return ErrInvalidCounterValue
		}
		return u.storage.UpdateCounter(m.ID, *m.Counter)
	case GaugeMetricType:
		if m.Gauge == nil {
			return ErrInvalidGaugeValue
		}
		return u.storage.UpdateGauge(m.ID, *m.Gauge)
	}
	return ErrUnknownMetricType
}
