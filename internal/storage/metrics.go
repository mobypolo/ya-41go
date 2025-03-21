package storage

import "github.com/mobypolo/ya-41go/internal/customerrors"

type GaugeMetric string
type CounterMetric string

const (
	GaugeTemperature GaugeMetric = "temperature"
	GaugeLoad        GaugeMetric = "load"

	CounterRequests CounterMetric = "requests"
	CounterErrors   CounterMetric = "errors"
)

func (g GaugeMetric) String() string {
	return string(g)
}
func (c CounterMetric) String() string {
	return string(c)
}

func ParseGaugeMetric(s string) (GaugeMetric, error) {
	switch s {
	case GaugeTemperature.String():
		return GaugeTemperature, nil
	case GaugeLoad.String():
		return GaugeLoad, nil
	default:
		return "", customerrors.ErrUnknownGaugeName
	}
}

func ParseCounterMetric(s string) (CounterMetric, error) {
	switch s {
	case CounterRequests.String():
		return CounterRequests, nil
	case CounterErrors.String():
		return CounterErrors, nil
	default:
		return "", customerrors.ErrUnknownCounterName
	}
}
