package agent

import (
	m "metrics/internal/models"
	"time"
)

type Worker struct {
	client    UserClient
	collector Collector
}

func NewWorker(config AgentConfig) Worker {
	client := NewUserClient(config)

	collector := NewCollector(config)

	return Worker{
		client:    client,
		collector: collector,
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
			w.client.SendMetricContainer(container)
		case <-pollTicker.C:
			w.collector.CollectMetrics(&pollCount, &container)
		}
	}

}
