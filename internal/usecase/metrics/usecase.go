package metrics

import "strconv"

type Repository interface {
	UpdateCounter(name string, value Counter) error
	UpdateGuage(name string, value Gauge) error
}

type Usecase struct {
	repo Repository
}

func NewUsecase(repo Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Update(mtype, mname, mvalue string) error {
	if mname == "" {
		return ErrUnknownMetric
	}
	switch mtype {
	case CounterName:
		value, err := strconv.ParseInt(mvalue, 10, 64)
		if err != nil {
			return ErrBadMetricValue
		}
		return u.repo.UpdateCounter(mname, Counter(value))
	case GaugeName:
		value, err := strconv.ParseFloat(mvalue, 64)
		if err != nil {
			return ErrBadMetricValue
		}
		return u.repo.UpdateGuage(mname, Gauge(value))
	}
	return ErrUnknownMetricType
}
