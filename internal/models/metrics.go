package models

import (
	"log"
	"strconv"
)

type MetricType string

const (
	CounterType MetricType = "counter"
	GaugeType   MetricType = "gauge"
)

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

// NewMetricSendContainer is used for collecting data before sending to server
type MetricSendContainer struct {
	GaugeMetrics   map[string]string
	CounterMetrics map[string]string
	UserMetrics    map[string]string
}

func (mc *MetricSendContainer) ConvertContainerToUpdateMetricsModel() []UpdateMetricsModel {

	updateMetrics := make([]UpdateMetricsModel, 0, len(mc.GaugeMetrics)+len(mc.CounterMetrics)+len(mc.UserMetrics))

	for metric, value := range mc.GaugeMetrics {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println("!Error while parsing float value ", err, " for metric : ", metric, " value: ", value)
			value = 0
		}

		updateMetrics = append(updateMetrics, UpdateMetricsModel{
			ID:    metric,
			MType: string(GaugeType),
			Value: &value,
		})
	}
	for metric, value := range mc.CounterMetrics {
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Println("!Error while parsing int value ", err, " for metric : ", metric)
			value = 0
		}
		updateMetrics = append(updateMetrics, UpdateMetricsModel{
			ID:    metric,
			MType: string(CounterType),
			Delta: &value,
		})
	}
	for metric, value := range mc.UserMetrics {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println("!Error while parsing float value ", err, " for metric : ", metric, " value: ", value)
			value = 0
		}

		updateMetrics = append(updateMetrics, UpdateMetricsModel{
			ID:    metric,
			MType: string(GaugeType),
			Value: &value,
		})
	}

	return updateMetrics
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

	userMap := make(map[string]string)
	userMap["RandomValue"] = ""

	metricContainer.UserMetrics = userMap
	metricContainer.GaugeMetrics = gaugeMap
	metricContainer.CounterMetrics = counterMap
	return metricContainer
}

// MetricsDict is used for db initialization
var MetricsDict map[MetricType][]string

func init() {
	MetricsDict = make(map[MetricType][]string)
	MetricsDict[CounterType] = []string{"PollCount"}
	MetricsDict[GaugeType] = append(MetricsDict[GaugeType], SystemMetrics...)
	MetricsDict[GaugeType] = append(MetricsDict[GaugeType], "RandomValue")

}
