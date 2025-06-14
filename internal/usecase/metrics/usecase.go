package metrics

type MetricsStorage interface {
	GetCounter(name string) ([]int64, error)
	GetGauge(name string) (float64, error)
	GetCounters() map[string][]int64
	GetGauges() map[string]float64
	AppendCounter(name string, value int64) error
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
	Counters map[string][]int64
	Gauges   map[string]float64
}

func (u *Usecase) GetMetrics() Metrics {
	return Metrics{
		Counters: u.storage.GetCounters(),
		Gauges:   u.storage.GetGauges(),
	}
}

// GetCounterSum returns the sum of values for the specified counter.
func (u *Usecase) GetCounterSum(name string) (int64, error) {
	values, err := u.storage.GetCounter(name)
	if err != nil {
		return 0, err
	}
	var sum int64
	for _, value := range values {
		sum += value
	}
	return sum, nil
}

// GetGauge returnes the value for the specified gauge.
func (u *Usecase) GetGauge(name string) (float64, error) {
	value, err := u.storage.GetGauge(name)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// AppendCounter appends a new value to the specified counter's value list.
func (u *Usecase) AppendCounter(name string, value int64) error {
	return u.storage.AppendCounter(name, value)
}

// UpdateGaute replaces the previous metric value with a new one.
func (u *Usecase) UpdateGauge(name string, value float64) error {
	return u.storage.UpdateGauge(name, value)
}
