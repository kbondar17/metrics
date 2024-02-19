package agent

import (
	m "metrics/internal/models"
	"time"

	"go.uber.org/zap"
)

type Worker struct {
	client    UserClient
	collector Collector
	logger    *zap.SugaredLogger
}

func NewWorker(config AgentConfig, logger *zap.SugaredLogger) Worker {
	client := NewUserClient(config)

	collector := NewCollector(config)

	return Worker{
		client:    client,
		collector: collector,
		logger:    logger,
	}
}

func (w Worker) Run() {
	var pollCount int
	container := m.NewMetricSendContainer()

	reportTicker := time.NewTicker(time.Duration(w.collector.config.reportInterval) * time.Second)
	defer reportTicker.Stop()

	pollTicker := time.NewTicker(time.Duration(w.collector.config.pollInterval) * time.Second)
	defer pollTicker.Stop()

	// w.collector.CollectMetrics(&pollCount, &container)
	// w.client.SendMetricContainerInButches(container)

	for {
		select {
		case <-reportTicker.C:
			//TODO: как соблюссти обратную совместимость?
			w.client.SendBoth(container, w.logger)
			// w.client.SendMetricContainer(container)
			// w.client.SendMetricContainerInButches(container, w.logger)
		case <-pollTicker.C:
			w.collector.CollectMetrics(&pollCount, &container, w.logger)
		}
	}

}
