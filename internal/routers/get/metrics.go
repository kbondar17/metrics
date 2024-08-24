package get

import (
	er "metrics/internal/errors"
	"metrics/internal/models"
	repo "metrics/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetCount(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType, logger *zap.SugaredLogger) {

	rg.GET("/:name", func(c *gin.Context) {
		metricName := c.Params.ByName("name")
		metric, err := repository.GetCountMetricValueByName(metricName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, metric)
	})

}

func GetGauge(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType, logger *zap.SugaredLogger) {

	rg.GET("/:name", func(c *gin.Context) {
		metricName := c.Params.ByName("name")
		metric, err := repository.GetGaugeMetricValueByName(metricName, metricType)
		if err == er.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"metric name": metricName})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(200, metric)
	})

}
