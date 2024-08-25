package agent

import (
	_ "encoding/json"
	"fmt"
	"math/rand"
	m "metrics/internal/models"
	"reflect"
	"runtime"
	"sync"
	"time"

	gopsutil "github.com/shirou/gopsutil/v4/mem"
	"go.uber.org/zap"
)

type Collector struct {
	config AgentConfig
	logger *zap.SugaredLogger
}

func NewCollector(config AgentConfig, logger *zap.SugaredLogger) Collector {
	return Collector{config: config, logger: logger}
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

func (coll *Collector) collectAdditionalMetrics(container *m.MetricSendContainer, wg *sync.WaitGroup) {
	v, _ := gopsutil.VirtualMemory()
	container.GaugeMetrics["TotalMemory"] = fmt.Sprintf("%d", v.Total)
	container.GaugeMetrics["FreeMemory"] = fmt.Sprintf("%d", v.Free)
	container.GaugeMetrics["CPUutilization1"] = fmt.Sprintf("%f", v.UsedPercent)
	// time.Sleep(10 * time.Second)
}

func (coll *Collector) CollectMetrics(pollCount *int32, pollInterval int, reportInterval int, dataChan chan<- m.MetricSendContainer) {
	container := m.NewMetricSendContainer()
	pollTicker := time.NewTicker(time.Duration(pollInterval) * time.Second)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer reportTicker.Stop()

	var mu sync.Mutex
	var mem runtime.MemStats
	wg := sync.WaitGroup{}
	for {
		select {
		case <-reportTicker.C:
			mu.Lock()
			dataChan <- container
			mu.Unlock()
		case <-pollTicker.C:
			newContainer := m.NewMetricSendContainer()
			for metric := range newContainer.GaugeMetrics {
				runtime.ReadMemStats(&mem)
				v := reflect.ValueOf(mem)
				metricValueRaw := v.FieldByName(metric)
				mu.Lock()
				newContainer.GaugeMetrics[metric] = parseMetric(metric, metricValueRaw, coll.logger)
				mu.Unlock()
			}
			mu.Lock()
			newContainer.UserMetrcs["RandomValue"] = fmt.Sprintf("%f", rand.Float64())
			*pollCount++
			newContainer.CounterMetrics["PollCount"] = fmt.Sprintf("%d", *pollCount)
			mu.Unlock()
			wg.Add(1)
			go func() {
				defer wg.Done()
				coll.collectAdditionalMetrics(&newContainer, &wg)
			}()
			wg.Wait()
			mu.Lock()
			container = newContainer
			mu.Unlock()
			coll.logger.Infoln("Metrics collected")

		}
	}
}
