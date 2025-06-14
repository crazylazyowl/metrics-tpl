package api

import (
	"errors"

	"github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
)

type MetricUpdateReq struct {
	ID         string   `json:"id"`
	MetricType string   `json:"type"`
	Delta      *int64   `json:"delta,omitempty"`
	Value      *float64 `json:"value,omitempty"`
}

func (m MetricUpdateReq) Validate() error {
	switch m.MetricType {
	case metrics.CounterMetricType:
		if m.Delta == nil {
			return errors.New("delta is missing")
		}
	case metrics.GaugeMetricType:
		if m.Value == nil {
			return errors.New("value is missing")
		}
	default:
		return metrics.ErrUnknownMetricType
	}
	return nil
}

type MetricGetReq struct {
	ID         string   `json:"id"`
	MetricType string   `json:"type"`
	Delta      *int64   `json:"delta,omitempty"`
	Value      *float64 `json:"value,omitempty"`
}

func (m MetricGetReq) Validate() error {
	switch m.MetricType {
	case metrics.CounterMetricType, metrics.GaugeMetricType:
	default:
		return metrics.ErrUnknownMetricType
	}
	return nil
}
