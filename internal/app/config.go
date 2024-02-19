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
	DBDNS            string
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
	host := flag.String("a", defaultHost, "HTTP server address. Default is localhost:8080")
	if ennvHost := os.Getenv("ADDRESS"); ennvHost != "" {
		host = &ennvHost
	}

	defaultStoreInteval := 300

	storeInterval := flag.Int("i", defaultStoreInteval, "Interval for saving metrics to the database. Default is 300 seconds")
	if ennvStoreInterval := os.Getenv("STORE_INTERVAL"); ennvStoreInterval != "" {
		val, err := strconv.Atoi(ennvStoreInterval)
		if err != nil {
			panic(err)
		}
		storeInterval = &val
	}

	defaultStoragePath := "/tmp/metrics-db.json"

	storagePath := flag.String("f", defaultStoragePath, "Full file name where current values are saved. Default is /tmp/metrics-db.json, empty value disables disk writing function")
	if ennvStoragePath := os.Getenv("FILE_STORAGE_PATH"); ennvStoragePath != "" {
		storagePath = &ennvStoragePath
	}

	defaultRetorePolicy := true
	restoreOnStartUp := flag.Bool("r", defaultRetorePolicy, "Whether to restore the previous state from a file. Default is true")
	if ennvRestorePolicy := os.Getenv("RESTORE"); ennvRestorePolicy != "" {
		val, err := strconv.ParseBool(ennvRestorePolicy)
		if err != nil {
			panic(err)
		}
		restoreOnStartUp = &val
	}

	defaultDBDNS := ""
	dbDNS := flag.String("d", defaultDBDNS, "Database dns. Default is empty value.")

	if envDBDNS := os.Getenv("DATABASE_DSN"); envDBDNS != "" {
		dbDNS = &envDBDNS
	}

	flag.Parse()

	var mustSync bool

	if *storeInterval == 0 {
		mustSync = true
	} else {
		mustSync = false
	}

	storageConf := StorageConf{StoragePath: *storagePath, RestoreOnStartUp: *restoreOnStartUp, MustSync: mustSync, StoreInterval: *storeInterval, DBDNS: *dbDNS}
	return NewAppConfig(*host, storageConf)
}
