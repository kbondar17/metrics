package agent

import (
	"net/url"
)

type AgentConfig struct {
	PollInterval   int
	ReportInterval int
	ServerAddress  url.URL
}

func NewAgentConfig(pollInterval int, reportInterval int, serverAddress string) AgentConfig {
	u, err := url.Parse(serverAddress)
	if err != nil {
		panic(err)
	}

	return AgentConfig{
		PollInterval:   pollInterval,
		ReportInterval: reportInterval,
		ServerAddress:  *u,
	}
}
