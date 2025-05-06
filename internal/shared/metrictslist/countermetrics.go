package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",

	"testSetGet*",
	"PopulateCounter*",
}

func AllCounterMetricNames() []string {
	return append([]string{}, AllowedCounterMetrics...)
}
