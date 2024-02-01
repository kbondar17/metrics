package app

import (
	"flag"
	"os"
	"strconv"
)

type AppConfig struct {
	host             string
	StoreInterval    int
	StoragePath      string
	restoreOnStartUp bool
}

func NewAppConfig(host string, storeInterval int, storagePath string, restoreOnStartUp bool) *AppConfig {
	return &AppConfig{host: host, StoreInterval: storeInterval, StoragePath: storagePath, restoreOnStartUp: restoreOnStartUp}

}

func getConfig() (string, int, string, bool) {
	defaultHost := "localhost:8080"
	host := flag.String("a", defaultHost, "Адрес HTTP-сервера. По умолчанию localhost:8080")
	if ennvHost := os.Getenv("ADDRESS"); ennvHost != "" {
		host = &ennvHost
	}

	defaultStoreInteval := 300

	storeInterval := flag.Int("i", defaultStoreInteval, "Интервал сохранения метрик в БД. По умолчанию 300 секунд")
	if ennvStoreInterval := os.Getenv("STORE_INTERVAL"); ennvStoreInterval != "" {
		val, err := strconv.Atoi(ennvStoreInterval)
		if err != nil {
			panic(err)
		}
		storeInterval = &val
	}

	defaultStoragePath := "/tmp/metrics-db.json"
	// defaultStoragePath := "/Users/makbuk/go/src/yandex/metrics/internal/database/backup.json"

	storagePath := flag.String("f", defaultStoragePath, "Полное имя файла, куда сохраняются текущие значения. По умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск")
	if ennvStoragePath := os.Getenv("FILE_STORAGE_PATH"); ennvStoragePath != "" {
		storagePath = &ennvStoragePath
	}

	defaultRetorePolicy := true
	restoreOnStartUp := flag.Bool("r", defaultRetorePolicy, "Восстанавливать ли предыдущее состояние из файла. По умолчанию true")
	if ennvRestorePolicy := os.Getenv("RESTORE"); ennvRestorePolicy != "" {
		val, err := strconv.ParseBool(ennvRestorePolicy)
		if err != nil {
			panic(err)
		}
		restoreOnStartUp = &val
	}
	flag.Parse()
	return *host, *storeInterval, *storagePath, *restoreOnStartUp
}

func NewAppConfigFromEnv() *AppConfig {
	host, storeInterval, storagePath, restoreOnStartUp := getConfig()
	return NewAppConfig(host, storeInterval, storagePath, restoreOnStartUp)
}
