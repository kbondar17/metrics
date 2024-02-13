package app

import (
	"errors"
	"fmt"
	"log"
	"time"

	db "metrics/internal/database"
	"metrics/internal/database/memory"
	postgres "metrics/internal/database/postgres"
	er "metrics/internal/errors"
	logger "metrics/internal/logger"
	m "metrics/internal/models"
	repo "metrics/internal/repository"
	routes "metrics/internal/routers"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type App struct {
	Config     *AppConfig
	Router     *gin.Engine
	repository repo.MetricsCRUDer
	logger     *zap.SugaredLogger
}

func (a *App) Run() {
	a.logger.Info("Starting server on ", a.Config.host)
	a.Router.Run(a.Config.host)
}

// SaveDataInInterval saves data in interval
func (a *App) SaveDataInInterval(storeInterval int, fname string) {
	for {
		metrics, err := a.repository.GetAllMetrics()
		if err != nil {
			a.logger.Infof("failed to get metrics: %v", err)
		}
		for _, metric := range metrics {
			err := db.SaveMetric(fname, metric)
			if err != nil {
				a.logger.Infof("failed to save metric: %v", err)
			}
		}
		time.Sleep(time.Duration(storeInterval) * time.Second)
	}
}

// addDefaultMetrics creates all metrics in DB
func addDefaultMetrics(repository repo.MetricsCRUDer, logger *zap.SugaredLogger) {

	for metricType, metricArray := range m.MetricsDict {
		for _, name := range metricArray {
			err := repository.Create(name, metricType, logger)
			if err != nil && !errors.Is(err, er.ErrAlreadyExists) {
				logger.Infow("failed to create metric: %v", err)
			}
		}
	}

}

func NewApp(conf *AppConfig) *App {
	fmt.Println("!!!APP_CONFIG::", conf)

	logger, err := logger.NewAppLogger()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	var storage repo.Storager
	if conf.StorageConfig.DBDNS == "" {
		storage = memory.NewMemStorage()
		logger.Info("using memory storage")
	} else {
		storage, err = postgres.NewPostgresStorage(conf.StorageConfig.DBDNS)
		if err != nil {
			logger.Infoln("failed to create storage: %v", err)
		} else {
			logger.Info("using postgres storage")
		}
	}

	repository := repo.NewMerticsRepo(storage)

	if conf.StorageConfig.RestoreOnStartUp {
		restoredMetrics, err := db.Load(conf.StorageConfig.StoragePath)
		if err != nil {
			logger.Infof("failed to load metrics: %v", err)
		}
		logger.Infof("restored metrics: %v", restoredMetrics)
		for _, metric := range restoredMetrics {
			if metric.MType == string(m.GaugeType) {
				err := repository.UpdateMetric(metric.ID, m.GaugeType, *metric.Value, false, "", logger)
				if err != nil {
					logger.Infof("failed to update metric: %v", err)
				}
			}
			if metric.MType == string(m.CounterType) {
				err := repository.UpdateMetric(metric.ID, m.CounterType, *metric.Delta, false, "", logger)
				if err != nil {
					logger.Infof("failed to update metric: %v", err)
				}
			}
		}
		logger.Info("restored metrics")
	}

	router := routes.RegisterMerticsRoutes(repository, logger, conf.StorageConfig.MustSync, conf.StorageConfig.StoragePath)

	addDefaultMetrics(repository, logger)

	return &App{Config: conf, Router: router, repository: repository, logger: logger}
}
