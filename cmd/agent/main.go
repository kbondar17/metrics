package main

import (
	"flag"
	"fmt"
	"metrics/internal/agent"
	"metrics/internal/utils"
	"os"
	"strconv"
)

func getConfig() (int, int, string) {

	defaultHost := "http://localhost:8080"
	if host, exists := os.LookupEnv("HOST"); exists {
		defaultHost = host
	}

	host := flag.String("a", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")

	defaulreportInterval := 10
	if reportEnv, exists := os.LookupEnv("REPORT_INTERVAL"); exists {
		if reportInt, err := strconv.Atoi(reportEnv); err == nil {
			defaulreportInterval = reportInt
		}
	}

	reportInterval := flag.Int("r", defaulreportInterval, "Частота отправки метрик на сервер в секундах. По умолчанию 10")

	defaultPollInterval := 2
	if pollEnv, exists := os.LookupEnv("POLL_INTERVAL"); exists {
		if pollInt, err := strconv.Atoi(pollEnv); err == nil {
			defaultPollInterval = pollInt
		}
	}

	pollInterval := flag.Int("p", defaultPollInterval, "Частота опроса метрик в секундах. По умолчанию 2")

	flag.Parse()
	return *reportInterval, *pollInterval, *host

}

func main() {

	reportInterval, pollInterval, serverAddress := getConfig()

	fmt.Println(reportInterval, pollInterval, serverAddress)

	logger := utils.NewLogger("./logs/agent_logs.log", "Agent: ")

	agentConfig := agent.NewAgentConfig(pollInterval, reportInterval, serverAddress)
	worker := agent.NewWorker(agentConfig, logger)
	worker.Start()

}
