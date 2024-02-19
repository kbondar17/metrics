package agent

import (
	_ "encoding/json"
	"fmt"
	"math/rand"
	m "metrics/internal/models"
	"reflect"
	"runtime"

	"go.uber.org/zap"
)

type Collector struct {
	config AgentConfig
}

func NewCollector(config AgentConfig) Collector {
	return Collector{config: config}
}

func parseMetric(metricName string, value reflect.Value, logger *zap.SugaredLogger) string {
	if value.Type() == reflect.TypeOf(float64(0)) {
		return fmt.Sprintf("%f", value.Float())
	} else if value.Type() == reflect.TypeOf(uint64(0)) || value.Type() == reflect.TypeOf(uint32(0)) {
		return fmt.Sprintf("%d", value.Uint())
	} else {
		logger.Infof("Warning: Type of metric is neither int nor float: %s", metricName)
		return ""
	}
}

func (coll *Collector) CollectMetrics(pollCount *int, container *m.MetricSendContainer, logger *zap.SugaredLogger) {
	var mem runtime.MemStats

	for metric := range container.GaugeMetrics {
		runtime.ReadMemStats(&mem)
		v := reflect.ValueOf(mem)
		metricValueRaw := v.FieldByName(metric)
		container.GaugeMetrics[metric] = parseMetric(metric, metricValueRaw, logger)
	}

	container.UserMetrcs["RandomValue"] = fmt.Sprintf("%f", rand.Float64())

	*pollCount++

	container.CounterMetrics["PollCount"] = fmt.Sprintf("%d", *pollCount)

}
