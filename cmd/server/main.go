package main

import (
	"metrics/internal/app"
)

func getConfig() string {
	// defaultHost := "localhost:8080"
	// if host, exists := os.LookupEnv("HOST"); exists {
	// 	defaultHost = host
	// }

	// host := flag.String("host", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	// flag.Parse()
	// return *host
	return "localhost:8080"
}

func main() {

	host := getConfig()
	appConfig := app.NewAppConfig(host)
	app := app.NewApp(appConfig)

	app.Run()

}
