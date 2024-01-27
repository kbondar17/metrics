package app

import (
	"errors"
	"log"

	db "metrics/internal/database"
	er "metrics/internal/errors"
	m "metrics/internal/models"
	repo "metrics/internal/repository"
	routes "metrics/internal/routers"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config *AppConfig
	Router *gin.Engine
}

func (a *App) Run() {
	log.Println("Starting server on ", a.Config.host)
	a.Router.Run(a.Config.host)
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
	addDefaultMetrics(repository)
	return &App{Config: conf, Router: routes.RegisterMerticsRoutes(repository)}
}
