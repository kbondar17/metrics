package app

import (
	"errors"
	"log"
	"time"

	db "metrics/internal/database"
	er "metrics/internal/errors"
	logger "metrics/internal/logger"
	m "metrics/internal/models"
	repo "metrics/internal/repository"
	routes "metrics/internal/routers"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config     *AppConfig
	Router     *gin.Engine
	logger     *logger.AppLogger
	repository repo.MetricsCRUDer
}

func (a *App) Run() {
	log.Println("Starting server on ", a.Config.host)
	a.Router.Run(a.Config.host)
}

// SaveDataInInterval saves data in interval
func (a *App) SaveDataInInterval(storeInterval int, fname string) {
	for {
		metrics := a.repository.GetAllMetrics()

		for _, metric := range metrics {
			err := db.SaveMetric(fname, metric)
			if err != nil {
				log.Printf("failed to save metric: %v", err)
			}
		}
		time.Sleep(time.Duration(storeInterval) * time.Second)
	}
}

// addDefaultMetrics creates all metrics in DB
func addDefaultMetrics(repository repo.MetricsCRUDer) {

	for metricType, metricArray := range m.MetricsDict {
		for _, name := range metricArray {
			err := repository.Create(name, metricType)
			if err != nil && !errors.Is(err, er.ErrAlreadyExists) {
				log.Fatalf("failed to create metric: %v", err)
			}
		}
	}

}

func NewApp(conf *AppConfig) *App {
	storage := db.NewStorage()
	repository := repo.NewMerticsRepo(storage)
	logger := logger.NewAppLogger()

	if conf.StorageConfig.RestoreOnStartUp {
		restoredMetrics, err := db.Load(conf.StorageConfig.StoragePath)
		if err != nil {
			log.Printf("failed to load metrics: %v", err)
		}
		log.Println("restored metrics::", restoredMetrics)

		for _, metric := range restoredMetrics {
			if metric.MType == string(m.GaugeType) {
				err := repository.UpdateMetric(metric.ID, m.GaugeType, *metric.Value, false, "")
				if err != nil {
					log.Printf("failed to update metric: %v", err)
				}
			}
			if metric.MType == string(m.CounterType) {
				err := repository.UpdateMetric(metric.ID, m.CounterType, *metric.Delta, false, "")
				if err != nil {
					log.Printf("failed to update metric: %v", err)
				}
			}
		}
		log.Println("restored metrics")
	}

	router := routes.RegisterMerticsRoutes(repository, logger, conf.StorageConfig.MustSync, conf.StorageConfig.StoragePath)

	addDefaultMetrics(repository)

	return &App{Config: conf, Router: router, logger: logger, repository: repository}
}
