package agent

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

type Metric struct {
	Name  string
	Type  MetricType
	Value interface{} // float64 или int64
}
