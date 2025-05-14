package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",

	"testSetGet*",
	"PopulateCounter*",
	"CounterBatchZip*",
}

func AllCounterMetricNames() []string {
	return append([]string{}, AllowedCounterMetrics...)
}
