package app

import (
	"context"
	"log"

	db "metrics/internal/database"
	m "metrics/internal/models"
	repo "metrics/internal/repository"
	routes "metrics/internal/routers"
	"metrics/internal/utils"

	"github.com/gin-gonic/gin"
)

type App struct {
	Config *AppConfig
	Router *gin.Engine
}

func (a *App) Run() {
	log.Println("Starting server on ", a.Config.Server.Address)
	a.Router.Run(a.Config.Server.Address)
}

// addDefaultMetrics creates all metrics in DB
func addDefaultMetrics(repository repo.MetricsCRUDer) {

	for metricType, metricArray := range m.MetricsDict {
		for _, name := range metricArray {
			err := repository.Create(name, metricType)
			if err != nil && err != utils.AlreadyExists {
				log.Fatalf("failed to create metric: %v", err)
			}
		}
	}

}

func NewApp(conf *AppConfig) *App {
	ctx := context.Background()

	storage := db.NewStorage(ctx)
	repository := repo.NewMerticsRepo(storage)
	addDefaultMetrics(repository)
	return &App{Config: conf, Router: routes.RegisterMerticsRoutes(repository)}
}
