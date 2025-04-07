package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",

	"testSetGet*",
}

func AllCounterMetricNames() []string {
	return append([]string{}, AllowedCounterMetrics...)
}
