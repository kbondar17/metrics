package models

type MetricType string

const (
	CounterType MetricType = "counter"
	GaugeType   MetricType = "gauge"
)

type GaugeMetric struct {
	Name  string  `json:"name" validate:"required"`
	Type  string  `json:"type" validate:"required"`
	Value float64 `json:"value" validate:"required"`
}

type CounterMetric struct {
	Name  string `json:"name" validate:"required"`
	Type  string `json:"type" validate:"required"`
	Value int64  `json:"value" validate:"required"`
}

type Metric struct {
	Name  string     `json:"name" validate:"required"`
	Type  MetricType `json:"type" validate:"required"`
	Value float64    `json:"value" validate:"required"`
}

type MetricResponseModel struct {
	Name  string     `json:"name" validate:"required"`
	Type  MetricType `json:"type" validate:"required"`
	Value string     `json:"value" validate:"required"`
}

type MetricSendContainer struct {
	GaugeMetrics   map[string]string
	CounterMetrics map[string]string
	UserMetrcs     map[string]string
}

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
	userMap["RandomValue"] = ""

	metricContainer.UserMetrcs = userMap
	metricContainer.GaugeMetrics = gaugeMap
	metricContainer.CounterMetrics = counterMap
	return metricContainer

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

var MetricsDict map[MetricType][]string

func init() {
	MetricsDict = make(map[MetricType][]string)
	MetricsDict[CounterType] = []string{"PollCount"}
	MetricsDict[GaugeType] = []string{"RandomValue"}
	MetricsDict[GaugeType] = []string{"testGauge"}
	MetricsDict[GaugeType] = append(MetricsDict[GaugeType], SystemMetrics...)
}
