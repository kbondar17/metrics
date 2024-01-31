package agent

import (
	"fmt"
	"log"
	"math/rand"
	m "metrics/internal/models"
	"reflect"
	"runtime"
)

type Collector struct {
	config AgentConfig
}

func NewCollector(config AgentConfig) Collector {
	return Collector{config: config}
}

func parseMetric(metricName string, value reflect.Value) string {
	if value.Type() == reflect.TypeOf(float64(0)) {
		return fmt.Sprintf("%f", value.Float())
	} else if value.Type() == reflect.TypeOf(uint64(0)) || value.Type() == reflect.TypeOf(uint32(0)) {
		return fmt.Sprintf("%d", value.Uint())
	} else {
		log.Println("Warning: Type of metric is neither int nor float:", metricName)
		return ""
	}
}

func (coll *Collector) CollectMetrics(pollCount *int, container *m.MetricSendContainer) {
	var mem runtime.MemStats

	for metric := range container.GaugeMetrics {
		runtime.ReadMemStats(&mem)
		v := reflect.ValueOf(mem)
		metricValueRaw := v.FieldByName(metric)
		container.GaugeMetrics[metric] = parseMetric(metric, metricValueRaw)
	}

	container.UserMetrcs["RandomValue"] = fmt.Sprintf("%f", rand.Float64())

	*pollCount++

	container.CounterMetrics["PollCount"] = fmt.Sprintf("%d", *pollCount)
	log.Println("Collected data: ", container)

}
