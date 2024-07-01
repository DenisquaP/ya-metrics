package yametrics

// Metrics interface
type Metric interface {
	MetricGetter
	MetricWriter
	MetricSaver
}

// Iterface for writing metrics
type MetricWriter interface {
	WriteGauge(name string, val float64) error
	WriteCounter(name string, val int64) error
}

// Interface for getting metrics
type MetricGetter interface {
	GetMetrics() string
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

// Interface for saving metrics
type MetricSaver interface {
	SaveMetricsToFile(wd string) error
	RestoreFromFile(wd string) error
}
