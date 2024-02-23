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
	serverAddress  string
	sendStrategy   sendStrategy
}

type sendStrategy int

const (
	Single sendStrategy = iota
	Butches
)

func (s sendStrategy) String() string {
	return [...]string{"Single", "Butches"}[s]
}

func newAgentConfig(pollInterval int, reportInterval int, serverAddress string, dbDNS string) AgentConfig {
	u, err := url.Parse(serverAddress)
	if err != nil {
		panic(err)
	}

	var strategy sendStrategy
	if dbDNS != "" {
		strategy = Butches
	} else {
		strategy = Single
	}

	return AgentConfig{
		pollInterval:   pollInterval,
		reportInterval: reportInterval,
		serverAddress:  u.String(),
		sendStrategy:   strategy,
	}
}

func NewAgentConfigFromEnv() AgentConfig {
	reportInterval, pollInterval, serverAddress, dbDNS := parseConfig()
	log.Printf("Agent config: reportInterval: %d, pollInterval: %d, serverAddress: %s \n, dbDNS: %s", reportInterval, pollInterval, serverAddress, dbDNS)
	return newAgentConfig(pollInterval, reportInterval, serverAddress, dbDNS)
}

func parseConfig() (int, int, string, string) {
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

	httpHost := "http://" + *host

	defaultDBDNS := ""
	dbDNS := flag.String("d", defaultDBDNS, "Database dns. Default is empty value.")

	if envDBDNS := os.Getenv("DATABASE_DSN"); envDBDNS != "" {
		dbDNS = &envDBDNS
	}
	flag.Parse()

	return *reportInterval, *pollInterval, httpHost, *dbDNS
}
