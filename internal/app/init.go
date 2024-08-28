package app

import (
	"errors"
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
			a.logger.Errorf("failed to get metrics: %w", err)
		}
		for _, metric := range metrics {
			err := db.SaveMetric(fname, metric)
			if err != nil {
				a.logger.Errorf("failed to save metric: %w", err)
			}
		}
		time.Sleep(time.Duration(storeInterval) * time.Second)
	}
}

// addDefaultMetrics creates all metrics in DB
func addDefaultMetrics(repository repo.MetricsCRUDer, logger *zap.SugaredLogger) {

	for metricType, metricArray := range m.MetricsDict {
		for _, name := range metricArray {
			err := repository.Create(name, metricType)
			if err != nil && !errors.Is(err, er.ErrAlreadyExists) {
				logger.Errorf("failed to create metric: %w", err)
			}
		}
	}

}

func NewApp(conf *AppConfig) *App {
	logger, err := logger.New()

	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}

	logger.Infof("config: %v", conf)

	var storage repo.Storager
	switch {
	case conf.StorageConfig.DBDNS == "":
		storage = memory.NewMemStorage()
		logger.Info("using memory storage")
	default:
		storage, err = postgres.NewPostgresStorage(conf.StorageConfig.DBDNS, logger)
		if err != nil {
			logger.Fatalln("failed to create storage: %v", err)
		} else {
			logger.Info("using postgres storage")
		}
	}

	repository := repo.NewMerticsRepo(storage, logger)

	if conf.StorageConfig.RestoreOnStartUp {
		restoredMetrics, err := db.Load(conf.StorageConfig.StoragePath, logger)
		if err != nil {
			// не делаю return и не падаю с ошибкой, на случай, если файла не существует и загружать нечего
			logger.Infof("failed to load metrics: %w", err)
		}
		for _, metric := range restoredMetrics {
			if metric.MType == string(m.GaugeType) {
				if metric.Value != nil {
					err := repository.UpdateMetric(metric.ID, m.GaugeType, *metric.Value, false, "")
					if err != nil {
						logger.Infof("failed to update metric: %v", err)
					}
				} else {
					logger.Infof("metric.Value is nil for metric: %v", metric)
				}
			}
			if metric.MType == string(m.CounterType) {
				if metric.Delta != nil {
					err := repository.UpdateMetric(metric.ID, m.CounterType, *metric.Delta, false, "")
					if err != nil {
						logger.Infof("failed to update metric: %v", err)
					}
				} else {
					logger.Infof("metric.Delta is nil for metric: %v", metric)
				}
			}
		}
		logger.Info("restored metrics")

	}

	router := routes.RegisterMerticsRoutes(repository, logger, conf.StorageConfig.MustSync, conf.StorageConfig.StoragePath, conf.hashKey)

	addDefaultMetrics(repository, logger)

	return &App{Config: conf, Router: router, repository: repository, logger: logger}
}
