package main

import (
	"log"
	"metrics/internal/agent"
	"metrics/internal/logger"
)

func main() {
	config := agent.NewAgentConfigFromEnv()
	logger, err := logger.New()
	if err != nil {
		log.Fatal("error while creating logger ", err)
	}

	worker := agent.NewWorker(config, logger)
	worker.Run()
}
