package storage

import (
	"github.com/mobypolo/ya-41go/internal/server/customerrors"
	"github.com/mobypolo/ya-41go/internal/shared/metrictslist"
	"github.com/mobypolo/ya-41go/internal/shared/utils"
)

type GaugeMetric string
type CounterMetric string

func (g GaugeMetric) String() string {
	return string(g)
}
func (c CounterMetric) String() string {
	return string(c)
}

func ParseGaugeMetric(s string) (string, error) {
	exist := utils.ExistInArrayWithWildCard(metrictslist.AllGaugeMetricNames(), s)
	if exist {
		return s, nil
	}
	return "", customerrors.ErrUnknownGaugeName
}

func ParseCounterMetric(s string) (string, error) {
	exist := utils.ExistInArrayWithWildCard(metrictslist.AllCounterMetricNames(), s)
	if exist {
		return s, nil
	}
	return "", customerrors.ErrUnknownCounterName
}
