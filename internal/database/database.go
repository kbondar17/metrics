package database

type Metric struct {
	name string
}

type MetricLog struct {
	metric    Metric
	value     float64
	timestamp int64
}

type MemStorage struct {
	data map[string][]MetricLog
}
