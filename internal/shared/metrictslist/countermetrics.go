package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",
}

func AllCounterMetricNames() []string {
	keys := make([]string, 0, len(AllowedCounterMetrics))
	for i := range AllowedCounterMetrics {
		keys = append(keys, AllowedCounterMetrics[i])
	}
	return keys
}
