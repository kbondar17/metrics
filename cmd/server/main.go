package main

import (
	"metrics/internal/app"
)

func main() {
	appConfig := app.NewAppConfigFromEnv()
	app := app.NewApp(appConfig)

	app.Run()
}
