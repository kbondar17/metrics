package main

import (
	"flag"
	"fmt"
	"metrics/internal/app"
	"os"
)

func getConfig() string {
	defaultHost := "127.0.0.1:8080"
	if host, exists := os.LookupEnv("HOST"); exists {
		defaultHost = host
	}

	host := flag.String("host", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	flag.Parse()
	fmt.Println("Сервер запущен на", *host)
	return *host
}

func main() {

	host := getConfig()
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
