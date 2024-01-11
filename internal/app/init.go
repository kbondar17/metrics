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

	for _, name := range m.SystemMetrics {
		err := repository.Create(name, m.GaugeType)
		if err != nil && err != utils.AlreadyExists {
			log.Fatalf("failed to create metric: %v", err)
		}

	}
	err := repository.Create("PollCount", m.CounterType)

	if err != nil && err != utils.AlreadyExists {
		log.Fatalf("failed to create metric: %v", err)
	}

	err = repository.Create("testCounter", m.CounterType)

	if err != nil && err != utils.AlreadyExists {
		log.Fatalf("failed to create metric: %v", err)
	}

	err = repository.Create("RandomValue", m.GaugeType)

	if err != nil && err != utils.AlreadyExists {
		log.Fatalf("failed to create metric: %v", err)
	}

}

func NewApp(conf *AppConfig) *App {
	ctx := context.Background()

	storage := db.NewStorage(ctx)
	repository := repo.NewMerticsRepo(storage)
	addDefaultMetrics(repository)
	return &App{Config: conf, Router: routes.RegisterMerticsRoutes(repository)}
}
