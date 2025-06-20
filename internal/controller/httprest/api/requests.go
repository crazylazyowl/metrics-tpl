package api

import (
	"errors"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

type MetricUpdateReq struct {
	ID         string   `json:"id"`
	MetricType string   `json:"type"`
	Counter    *int64   `json:"delta,omitempty"`
	Gauge      *float64 `json:"value,omitempty"`
}

func (m MetricUpdateReq) Validate() error {
	if m.ID == "" {
		return errors.New("id is missing")
	}
	switch m.MetricType {
	case metrics.CounterMetricType:
		if m.Counter == nil {
			return errors.New("counter value is missing")
		}
	case metrics.GaugeMetricType:
		if m.Gauge == nil {
			return errors.New("gauge value is missing")
		}
	default:
		return metrics.ErrUnknownMetricType
	}
	return nil
}

type MetricGetReq struct {
	ID         string   `json:"id"`
	MetricType string   `json:"type"`
	Counter    *int64   `json:"delta,omitempty"`
	Gauge      *float64 `json:"value,omitempty"`
}

func (m MetricGetReq) Validate() error {
	if m.ID == "" {
		return errors.New("id is missing")
	}
	switch m.MetricType {
	case metrics.CounterMetricType:
	case metrics.GaugeMetricType:
	default:
		return metrics.ErrUnknownMetricType
	}
	return nil
}
