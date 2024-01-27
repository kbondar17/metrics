package app

import (
	"flag"
	"os"
)

type AppConfig struct {
	host string
}

func NewAppConfig(host string) *AppConfig {
	return &AppConfig{host: host}

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
