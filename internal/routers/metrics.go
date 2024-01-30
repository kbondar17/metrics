package routers

import (
	"fmt"
	"log"
	er "metrics/internal/errors"

	logger "metrics/internal/logger"
	"metrics/internal/models"
	repo "metrics/internal/repository"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// registerUpdateCounterRoutes registers handlers for metrics of type 'Counter'
func registerUpdateCounterRoutes(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType) {

	rg.POST("/:name/:value", func(c *gin.Context) {
		name := c.Params.ByName("name")
		value, err := strconv.Atoi(c.Params.ByName("value"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		err = repository.UpdateMetric(name, models.CounterType, value)
		if err == er.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
			return
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

func registerUpdateGaugeRoutes(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType) {

	rg.POST("/:name/:value", func(c *gin.Context) {
		name := c.Params.ByName("name")
		value, err := strconv.ParseFloat(c.Params.ByName("value"), 64)

		if err != nil {
			log.Println("error parsing path params", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		err = repository.UpdateMetric(name, models.GaugeType, value)
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

func registerGetCountRoutes(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType) {

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

func registerGetGaugeRoutes(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType) {

	rg.GET("/:name", func(c *gin.Context) {
		metricName := c.Params.ByName("name")
		metric, err := repository.GetGaugeMetricValueByName(metricName, metricType)
		if err == er.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"metric name": metricName, "error": "metric not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(200, metric)
	})

}

func RegisterGetValueRoute(rg *gin.RouterGroup, repository repo.MetricsCRUDer) {
	rg.POST("/", func(c *gin.Context) {
		var metric models.UpdateMetricsModel
		c.BindJSON(&metric)

		if metric.MType == string(models.GaugeType) {
			value, err := repository.GetGaugeMetricValueByName(metric.ID, models.GaugeType)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})
			}
			c.Header("Content-Type", "application/json")
			c.JSON(200, value)
		} else if metric.MType == string(models.CounterType) {
			value, err := repository.GetCountMetricValueByName(metric.ID)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})

			}
			c.Header("Content-Type", "application/json")
			c.JSON(200, value)

		}
	})
}

func RegisterUpdateRoute(rg *gin.RouterGroup, repository repo.MetricsCRUDer) {

	rg.POST("/", func(c *gin.Context) {
		var metric models.UpdateMetricsModel

		// TODO: err handling
		c.BindJSON(&metric)

		if metric.MType == string(models.GaugeType) {
			err := repository.UpdateMetric(metric.ID, models.GaugeType, *metric.Value)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})
			}

		}
		// return body
		c.JSON(200, metric)

	})

}

func RegisterMerticsRoutes(repository repo.MetricsCRUDer, logger *logger.AppLogger) *gin.Engine {

	r := gin.New()
	r.Use(RequestLogger(logger))

	r.POST("/echo", func(c *gin.Context) {
		//parse body and send it back
		var body interface{}
		c.BindJSON(&body)
		body.(map[string]interface{})["message"] = "from server"
		c.JSON(200, body)

	})

	r.LoadHTMLFiles("templates/metrics.html")

	r.GET("/", func(c *gin.Context) {
		metrics := repository.GetAllMetrics()
		fmt.Println("metrics:", metrics)
		c.HTML(http.StatusOK, "metrics.html", gin.H{
			"metrics": metrics,
		})

	})

	updateGroup := r.Group("/update")
	RegisterUpdateRoute(updateGroup, repository)
	registerUpdateGaugeRoutes(updateGroup.Group("/gauge"), repository, models.GaugeType)
	registerUpdateCounterRoutes(updateGroup.Group("/counter"), repository, models.CounterType)

	getGroup := r.Group("/value")
	RegisterGetValueRoute(getGroup, repository)
	registerGetGaugeRoutes(getGroup.Group("/gauge"), repository, models.GaugeType)
	registerGetCountRoutes(getGroup.Group("/counter"), repository, models.CounterType)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not Found"})
	})

	return r
}

func guidMiddleware() gin.HandlerFunc {
	fmt.Println("guidMiddleware is called")

	return func(c *gin.Context) {
		uuid := uuid.New()
		c.Set("uuid", uuid)
		fmt.Printf("The request with uuid %s is started \n", uuid)
		c.Next()
		fmt.Printf("The request with uuid %s is served \n", uuid)
	}
}
