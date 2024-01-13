package main

import (
	"flag"
	"metrics/internal/app"
	"os"
)

func getConfig() string {
	defaultHost := "localhost:8080"
	if host, exists := os.LookupEnv("HOST"); exists {
		defaultHost = host
	}

	host := flag.String("a", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	flag.Parse()
	return *host
}

func main() {

	host := getConfig()
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
