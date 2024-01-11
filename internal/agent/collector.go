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

var PollCount int

type Collector struct {
	config AgentConfig
	logger *log.Logger
}

func NewCollector(logger *log.Logger, config AgentConfig) Collector {
	return Collector{logger: logger, config: config}
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
func (self *Collector) CollectMetrics() m.MetricSendContainer {
	var mem runtime.MemStats

	container := m.NewMetricSendContainer()
	start := time.Now()

	for {
		PollCount++

		for metric, _ := range container.GaugeMetrics {
			runtime.ReadMemStats(&mem)
			v := reflect.ValueOf(mem)
			metricValueRaw := v.FieldByName(metric)
			container.GaugeMetrics[metric] = parseMetric(metric, metricValueRaw)
		}

		container.UserMetrcs["RandomValue"] = fmt.Sprintf("%f", rand.Float64())

		if int(time.Since(start).Seconds()) >= self.config.ReportInterval {
			self.logger.Println("Time to send metrics. Collected data: ", container)
			start = time.Now()
			break
		}
		time.Sleep(time.Duration(self.config.PollInterval) * time.Second)

	}

	container.CounterMetrics["PollCount"] = fmt.Sprintf("%d", PollCount)
	return container
}
