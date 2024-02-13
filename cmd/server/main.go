package main

import (
	_ "encoding/json"
	"metrics/internal/app"
)

func main() {
	appConfig := app.NewAppConfigFromEnv()
	app := app.NewApp(appConfig)

	if app.Config.StorageConfig.StoreInterval > 0 {
		go func() {
			app.SaveDataInInterval(app.Config.StorageConfig.StoreInterval, app.Config.StorageConfig.StoragePath)
		}()
	}

	app.Run()
}
