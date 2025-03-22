package metrictslist

var AllowedGaugeMetrics = []string{
	"temperature",
	"load",
	"testGauge",

	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",

	"RandomValue",
}

func AllGaugeMetricNames() []string {
	keys := make([]string, 0, len(AllowedGaugeMetrics))
	for _, item := range AllowedGaugeMetrics {
		keys = append(keys, item)
	}
	return keys
}
