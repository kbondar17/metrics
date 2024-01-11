package agent

import (
	"log"
)

type Worker struct {
	client    UserClient
	collector Collector
}

func NewWorker(config AgentConfig, logger *log.Logger) Worker {
	client := NewUserClient(config, logger)

	collector := NewCollector(logger, config)

	return Worker{
		client:    client,
		collector: collector,
	}
}

func (w Worker) Start() {
	for {
		data := w.collector.CollectMetrics()
		w.client.SendMetricContainer(data)
	}

}
