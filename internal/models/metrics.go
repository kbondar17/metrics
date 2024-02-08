package models

type MetricType string

const (
	CounterType MetricType = "counter"
	GaugeType   MetricType = "gauge"
)

type GaugeMetric struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type CounterMetric struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value int64  `json:"value"`
}

type Metric struct {
	Name  string     `json:"name"`
	Type  MetricType `json:"type"`
	Value float64    `json:"value"`
}

type MetricResponseModel struct {
	Name  string     `json:"name"`
	Type  MetricType `json:"type"`
	Value string     `json:"value"`
}

type UpdateMetricsModel struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricSendContainer struct {
	GaugeMetrics   map[string]string
	CounterMetrics map[string]string
	UserMetrcs     map[string]string
}

var (
	SystemMetrics = []string{
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
	}

	AllMetricsNames = append(SystemMetrics, "RandomValue", "PollCount", "testCounter", "testGauge")
)

// NewMetricSendContainer is used for collecting data before sending to server
func NewMetricSendContainer() MetricSendContainer {
	var metricContainer MetricSendContainer

	gaugeMap := make(map[string]string)

	for _, metricName := range SystemMetrics {
		gaugeMap[metricName] = ""
	}
	counterMap := make(map[string]string)
	counterMap["PollCount"] = ""
	counterMap["testCounter"] = ""

	userMap := make(map[string]string)
	userMap["testGauge"] = ""
	userMap["RandomValue"] = ""

	metricContainer.UserMetrcs = userMap
	metricContainer.GaugeMetrics = gaugeMap
	metricContainer.CounterMetrics = counterMap
	return metricContainer

}

// MetricsDict is used for db initialization
var MetricsDict map[MetricType][]string

func init() {
	MetricsDict = make(map[MetricType][]string)
	MetricsDict[CounterType] = []string{"PollCount", "testCounter"}
	MetricsDict[GaugeType] = append(MetricsDict[GaugeType], SystemMetrics...)
	MetricsDict[GaugeType] = append(MetricsDict[GaugeType], "RandomValue", "testGauge")

}
