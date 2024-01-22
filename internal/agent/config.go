package agent

import (
	"flag"
	"log"
	"net/url"
	"os"
	"strconv"
)

type AgentConfig struct {
	pollInterval   int
	reportInterval int
	serverAddress  url.URL
}

func NewAgentConfig(pollInterval int, reportInterval int, serverAddress string) AgentConfig {
	u, err := url.Parse(serverAddress)
	if err != nil {
		panic(err)
	}

	return AgentConfig{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddress:  *u,
	}
}

func NewAgentConfigFromEnv() AgentConfig {
	reportInterval, pollInterval, serverAddress := parseConfig()
	log.Printf("Agent config: reportInterval: %d, pollInterval: %d, serverAddress: %s \n", reportInterval, pollInterval, serverAddress)
	return NewAgentConfig(pollInterval, reportInterval, serverAddress)
}

func parseConfig() (int, int, string) {

	defaultHost := "localhost:8080"
	if host, exists := os.LookupEnv("ADDRESS"); exists {
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

	httpHost := "http://" + *host

	return *reportInterval, *pollInterval, httpHost

}
