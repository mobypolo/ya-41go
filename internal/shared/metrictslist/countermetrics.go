package metrictslist

var AllowedCounterMetrics = []string{
	"requests",
	"errors",
	"testCounter",

	"PollCount",
}

func AllCounterMetricNames() []string {
	keys := make([]string, 0, len(AllowedCounterMetrics))
	for _, item := range AllowedCounterMetrics {
		keys = append(keys, item)
	}
	return keys
}
