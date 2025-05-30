package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",

	"testSetGet*",
	"PopulateCounter*",
	"CounterBatchZip*",
	"GetSetZip*",
}

func AllCounterMetricNames() []string {
	return append([]string{}, AllowedCounterMetrics...)
}
