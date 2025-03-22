package repositories

type MetricsRepository interface {
	UpdateGauge(name string, value float64) error
	UpdateCounter(name string, delta int64) error
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}
