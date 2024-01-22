package agent

import (
	"fmt"
	"log"
	"math/rand"
	m "metrics/internal/models"
	"reflect"
	"runtime"
	"time"
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

// CollectMetrics collects metrics from runtime and returns them in MetricSendContainer
func (coll *Collector) _CollectMetrics() m.MetricSendContainer {
	var mem runtime.MemStats
	var pollCount int

	container := m.NewMetricSendContainer()
	start := time.Now()

	for {
		pollCount++

		for metric := range container.GaugeMetrics {
			runtime.ReadMemStats(&mem)
			v := reflect.ValueOf(mem)
			metricValueRaw := v.FieldByName(metric)
			container.GaugeMetrics[metric] = parseMetric(metric, metricValueRaw)
		}

		container.UserMetrcs["RandomValue"] = fmt.Sprintf("%f", rand.Float64())

		if int(time.Since(start).Seconds()) >= coll.config.reportInterval {
			log.Println("Time to send metrics. Collected data: ", container)
			break
		}
		time.Sleep(time.Duration(coll.config.pollInterval) * time.Second)

	}

	container.CounterMetrics["PollCount"] = fmt.Sprintf("%d", pollCount)
	return container
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

	container.CounterMetrics["PollCount"] = fmt.Sprintf("%d", pollCount)
	log.Println("Collected data: ", container)

}
