package routers

import (
	"log"
	"metrics/internal/app_errors"
	"metrics/internal/models"
	repo "metrics/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		if err == app_errors.ErrorNotFound {
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
		if err == app_errors.ErrorNotFound {
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
		if err == app_errors.ErrorNotFound {
			c.JSON(http.StatusNotFound, gin.H{"metric name": metricName, "error": "metric not found"})
			return
		}
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.JSON(200, metric)
	})

}

func RegisterMerticsRoutes(repository repo.MetricsCRUDer) *gin.Engine {
	r := gin.Default()

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
		c.HTML(http.StatusOK, "metrics.html", gin.H{
			"metrics": metrics,
		})

	})

	updateGroup := r.Group("/update")
	registerUpdateGaugeRoutes(updateGroup.Group("/gauge"), repository, models.GaugeType)
	registerUpdateCounterRoutes(updateGroup.Group("/counter"), repository, models.CounterType)

	getGroup := r.Group("/value")
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
