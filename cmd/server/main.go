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
	flag.Parse()

	if hostEnv, exists := os.LookupEnv("ADDRESS"); exists {
		host = &hostEnv
	}

	fmt.Println("Адрес HTTP-сервера: ", *host)
	return *host
}

func main() {

	host := getConfig()
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
