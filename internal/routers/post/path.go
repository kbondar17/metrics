package post

import (
	"log"
	er "metrics/internal/errors"
	"metrics/internal/models"
	repo "metrics/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// registerUpdateCounterRoutes registers handlers for metrics of type 'Counter'
func UpdateCounter(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType, logger *zap.SugaredLogger) {

	rg.POST("/:name/:value", func(c *gin.Context) {
		name := c.Params.ByName("name")
		value, err := strconv.ParseInt(c.Params.ByName("value"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err = repository.UpdateMetric(name, models.CounterType, value, false, "", logger)
		if err == er.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
			return
		}
		if err != nil {
			logger.Infow("failed to update metric", "name", name, "value", value, "error", err)
			c.Status(http.StatusBadRequest)
		}
		c.Status(http.StatusOK)

	})

	rg.POST("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

}

func UpdateGauge(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType, logger *zap.SugaredLogger) {

	rg.POST("/:name/:value", func(c *gin.Context) {
		name := c.Params.ByName("name")
		value, err := strconv.ParseFloat(c.Params.ByName("value"), 64)

		if err != nil {
			log.Println("error parsing path params", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err = repository.UpdateMetric(name, models.GaugeType, value, false, "", logger)
		if err == er.ErrorNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"metric name": name, "error": "metric not found"})
		}
		if err != nil {
			log.Println("error updating metric :: ", err)
			c.Status(http.StatusBadRequest)
		}
		c.Status(http.StatusOK)
	})

	rg.POST("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

}
