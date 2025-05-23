package metrics

type MetricsStorage interface {
	GetCounter(name string) ([]Counter, error)
	GetGauge(name string) (Gauge, error)
	GetCounters() map[string][]Counter
	GetGauges() map[string]Gauge
	UpdateCounter(name string, value Counter) error
	UpdateGauge(name string, value Gauge) error
}

type Usecase struct {
	storage MetricsStorage
}

func New(repo MetricsStorage) *Usecase {
	return &Usecase{storage: repo}
}

func (u *Usecase) Metrics() (map[string][]Counter, map[string]Gauge) {
	return u.storage.GetCounters(), u.storage.GetGauges()
}

func (u *Usecase) CounterSum(name string) (Counter, error) {
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

func (u *Usecase) Gauge(name string) (Gauge, error) {
	value, err := u.storage.GetGauge(name)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (u *Usecase) UpdateCounter(name string, value Counter) error {
	return u.storage.UpdateCounter(name, value)
}

func (u *Usecase) UpdateGauge(name string, value Gauge) error {
	return u.storage.UpdateGauge(name, value)
}
