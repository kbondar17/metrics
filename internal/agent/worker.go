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
	client := NewUserClient(config, logger)

	collector := NewCollector(config, logger)

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

	for {
		select {
		case <-reportTicker.C:
			if w.collector.config.sendStrategy == Single {
				w.client.SendMetricContainer(container)
			} else {
				w.client.SendMetricContainerInButches(container)
			}
		case <-pollTicker.C:
			w.collector.CollectMetrics(&pollCount, &container)
		}
	}

}
