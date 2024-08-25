package agent

import (
	m "metrics/internal/models"
	"sync"

	"go.uber.org/zap"
)

type Worker struct {
	client    UserClient
	collector Collector
	logger    *zap.SugaredLogger
}

func NewWorker(config AgentConfig, logger *zap.SugaredLogger) Worker {
	client := NewUserClient(config, logger)

	collector := NewCollector(config, logger)

	return Worker{
		client:    client,
		collector: collector,
		logger:    logger,
	}
}

func (w Worker) Run() {

	var pollCount int32

	wg := sync.WaitGroup{}

	dataChan := make(chan m.MetricSendContainer, 10)

	go w.collector.CollectMetrics(&pollCount, w.collector.config.pollInterval, w.collector.config.reportInterval, dataChan)
	go w.client.SendMetricContainer(dataChan)
	wg.Add(2)
	wg.Wait()

}
