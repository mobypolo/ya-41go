package sources

import (
	"github.com/mobypolo/ya-41go/internal/shared/utils"
	"runtime"
	"sync"
	"time"

	"github.com/mobypolo/ya-41go/internal/agent"
)

var mu sync.Mutex

var memStats runtime.MemStats

func collectMemStats() {
	mu.Lock()
	defer mu.Unlock()
	runtime.ReadMemStats(&memStats)
}

func getMemStatFloat64(f func(m *runtime.MemStats) float64) float64 {
	mu.Lock()
	defer mu.Unlock()
	return f(&memStats)
}

func getMemStatUint64(f func(m *runtime.MemStats) uint64) float64 {
	mu.Lock()
	defer mu.Unlock()
	return float64(f(&memStats))
}

type runtimeMetric struct {
	name string
	f    func() float64
}

func (r runtimeMetric) Name() string {
	return r.name
}

func (r runtimeMetric) Type() agent.MetricType {
	return agent.GaugeType
}

func (r runtimeMetric) Collect() (interface{}, error) {
	return r.f(), nil
}

func init() {
	go func() {
		for {
			collectMemStats()
			time.Sleep(time.Second)
		}
	}()

	metrics := []runtimeMetric{
		{"Alloc", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.Alloc }) }},
		{"BuckHashSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.BuckHashSys }) }},
		{"Frees", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.Frees }) }},
		{"GCCPUFraction", func() float64 { return getMemStatFloat64(func(m *runtime.MemStats) float64 { return m.GCCPUFraction }) }},
		{"GCSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.GCSys }) }},
		{"HeapAlloc", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.HeapAlloc }) }},
		{"HeapIdle", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.HeapIdle }) }},
		{"HeapInuse", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.HeapInuse }) }},
		{"HeapObjects", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.HeapObjects }) }},
		{"HeapReleased", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.HeapReleased }) }},
		{"HeapSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.HeapSys }) }},
		{"LastGC", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.LastGC }) }},
		{"Lookups", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.Lookups }) }},
		{"MCacheInuse", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.MCacheInuse }) }},
		{"MCacheSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.MCacheSys }) }},
		{"MSpanInuse", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.MSpanInuse }) }},
		{"MSpanSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.MSpanSys }) }},
		{"Mallocs", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.Mallocs }) }},
		{"NextGC", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.NextGC }) }},
		{"NumForcedGC", func() float64 {
			return getMemStatUint64(func(m *runtime.MemStats) uint64 { return uint64(m.NumForcedGC) })
		}},
		{"NumGC", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return uint64(m.NumGC) }) }},
		{"OtherSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.OtherSys }) }},
		{"PauseTotalNs", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.PauseTotalNs }) }},
		{"StackInuse", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.StackInuse }) }},
		{"StackSys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.StackSys }) }},
		{"Sys", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.Sys }) }},
		{"TotalAlloc", func() float64 { return getMemStatUint64(func(m *runtime.MemStats) uint64 { return m.TotalAlloc }) }},

		{"RandomValue", func() float64 { return float64(utils.RandInt(0, 100)) }},
	}

	for _, metric := range metrics {
		agent.Register(metric)
	}
}
