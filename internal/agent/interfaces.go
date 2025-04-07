package agent

type MetricSource interface {
	Name() string
	Type() MetricType
	Collect() (interface{}, error)
}
