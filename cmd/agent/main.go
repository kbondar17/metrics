package main

import (
	"metrics/internal/agent"
)

func main() {
	config := agent.NewAgentConfigFromEnv()
	worker := agent.NewWorker(config)
	worker.Run()
}
