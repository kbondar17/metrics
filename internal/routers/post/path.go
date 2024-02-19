package post

import (
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

		err = repository.UpdateMetric(name, models.CounterType, value, false, "")
		if err == er.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
			return
		}
		if err != nil {
			logger.Infof("error updating metric: %v", err)
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
			logger.Infof("error parsing path params: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err = repository.UpdateMetric(name, models.GaugeType, value, false, "")
		if err == er.ErrorNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"metric name": name})
		}
		if err != nil {
			logger.Infof("error updating metric: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.Status(http.StatusOK)
	})

	rg.POST("/", func(c *gin.Context) {
		c.Status(http.StatusNotFound)
	})

}
