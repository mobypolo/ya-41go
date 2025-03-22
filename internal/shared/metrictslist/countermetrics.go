package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",
}

func AllCounterMetricNames() []string {
	return append([]string{}, AllowedCounterMetrics...)
}
