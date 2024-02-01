package main

import (
	"metrics/internal/app"
)

func main() {
	appConfig := app.NewAppConfigFromEnv()
	app := app.NewApp(appConfig)

	// если не синхронная запись, тогда запускаем фоновую задачу
	if app.Config.StoreInterval > 0 {
		go func() {
			app.SaveDataInInterval(app.Config.StoreInterval, app.Config.StoragePath)
		}()
	}

	app.Run()
}
