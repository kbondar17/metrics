package main

import (
	_ "encoding/json"
	"log"
	"metrics/internal/agent"
	"metrics/internal/logger"
)

func main() {
	config := agent.NewAgentConfigFromEnv()
	logger, err := logger.NewAppLogger()
	if err != nil {
		log.Println("error while creating logger", "error", err)
	}

	worker := agent.NewWorker(config, logger)
	worker.Run()
}
