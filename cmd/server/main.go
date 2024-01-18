package main

import (
	"flag"
	"fmt"
	"metrics/internal/app"
	"os"
)

func getConfig() string {
	defaultHost := "localhost:8080"
	host := flag.String("a", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	if ennvHost := os.Getenv("ADDRESS"); ennvHost != "" {
		host = &ennvHost
	}
	flag.Parse()
	return *host
}

func main() {

	host := getConfig()
	fmt.Println("host: ", host)
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
