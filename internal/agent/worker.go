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

func (self Worker) Start() {
	for {
		data := self.collector.CollectMetrics()
		self.client.SendMetricContainer(data)
	}

}
