package app

import (
	"flag"
	"os"
	"strconv"
)

type StorageConf struct {
	StoragePath      string
	RestoreOnStartUp bool
	MustSync         bool
	StoreInterval    int
}

type AppConfig struct {
	host          string
	StorageConfig StorageConf
}

func NewAppConfig(host string, StorageConfig StorageConf) *AppConfig {
	return &AppConfig{host: host, StorageConfig: StorageConfig}
}

func NewAppConfigFromEnv() *AppConfig {
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

	var mustSync bool

	if *storeInterval == 0 {
		mustSync = true
	} else {
		mustSync = false
	}

	storageConf := StorageConf{StoragePath: *storagePath, RestoreOnStartUp: *restoreOnStartUp, MustSync: mustSync, StoreInterval: *storeInterval}
	return NewAppConfig(*host, storageConf)
}
