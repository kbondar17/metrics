package app

import (
	"flag"
	"os"
)

type ServerConfig struct {
	Address string
}

type AppConfig struct {
	Server ServerConfig
}

func NewAppConfig(host string) *AppConfig {
	return &AppConfig{
		Server: ServerConfig{
			Address: host,
		},
	}
}

func getConfig() string {
	defaultHost := "localhost:8080"
	host := flag.String("a", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	if ennvHost := os.Getenv("ADDRESS"); ennvHost != "" {
		host = &ennvHost
	}
	flag.Parse()
	return *host
}

func NewAppConfigFromEnv() *AppConfig {
	host := getConfig()
	return NewAppConfig(host)
}
