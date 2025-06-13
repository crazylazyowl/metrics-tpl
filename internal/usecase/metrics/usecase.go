package metrics

type MetricsStorage interface {
	GetCounter(name string) ([]Counter, error)
	GetGauge(name string) (Gauge, error)
	GetCounters() map[string][]Counter
	GetGauges() map[string]Gauge
	AppendCounter(name string, value Counter) error
	UpdateGauge(name string, value Gauge) error
}

type Usecase struct {
	storage MetricsStorage
}

func New(repo MetricsStorage) *Usecase {
	return &Usecase{storage: repo}
}

type Metrics struct {
	Counters map[string][]Counter
	Gauges   map[string]Gauge
}

func (u *Usecase) GetMetrics() Metrics {
	return Metrics{
		Counters: u.storage.GetCounters(),
		Gauges:   u.storage.GetGauges(),
	}
}

// GetCounterSum returns the sum of values for the specified counter.
func (u *Usecase) GetCounterSum(name string) (Counter, error) {
	values, err := u.storage.GetCounter(name)
	if err != nil {
		return 0, err
	}
	var sum Counter
	for _, value := range values {
		sum += value
	}
	return sum, nil
}

// GetGauge returnes the value for the specified gauge.
func (u *Usecase) GetGauge(name string) (Gauge, error) {
	value, err := u.storage.GetGauge(name)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// AppendCounter appends a new value to the specified counter's value list.
func (u *Usecase) AppendCounter(name string, value Counter) error {
	return u.storage.AppendCounter(name, value)
}

// UpdateGaute replaces the previous metric value with a new one.
func (u *Usecase) UpdateGauge(name string, value Gauge) error {
	return u.storage.UpdateGauge(name, value)
}
