package sources

import (
	"github.com/mobypolo/ya-41go/internal/agent"
	"github.com/mobypolo/ya-41go/internal/server/storage"
)

type counterMetric struct {
	name string
	f    func() float64
}

var PollCount = 1

func (r counterMetric) Name() string {
	return r.name
}

func (r counterMetric) Type() storage.MetricType {
	return storage.CounterType
}

func (r counterMetric) Collect() (interface{}, error) {
	return r.f(), nil
}

func init() {
	metrics := []counterMetric{
		{"PollCount", func() float64 { return float64(PollCount) }},
	}

	for _, metric := range metrics {
		agent.Register(metric)
	}
}
