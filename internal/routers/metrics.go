package routers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	er "metrics/internal/errors"
	"strings"

	logger "metrics/internal/logger"
	"metrics/internal/models"
	repo "metrics/internal/repository"

	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// registerUpdateCounterRoutes registers handlers for metrics of type 'Counter'
func registerUpdateCounterRoutes(rg *gin.RouterGroup, repository repo.MetricsCRUDer, metricType models.MetricType) {

	rg.POST("/:name/:value", func(c *gin.Context) {
		name := c.Params.ByName("name")
		value, err := strconv.ParseInt(c.Params.ByName("value"), 10, 64)

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
		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		}

		// десериализуем JSON
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse	body"})
			return
		}

		if metric.MType == string(models.GaugeType) {
			value, err := repository.GetGaugeMetricValueByName(metric.ID, models.GaugeType)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusNotFound, gin.H{"metric name": metric.ID, "error": "metric not found"})
			}
			c.Header("Content-Type", "application/json")
			c.JSON(200, models.UpdateMetricsModel{ID: metric.ID, MType: metric.MType, Value: &value})
		} else if metric.MType == string(models.CounterType) {
			value, err := repository.GetCountMetricValueByName(metric.ID)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusNotFound, gin.H{"metric name": metric.ID, "error": "metric not found"})

			}
			value64 := int64(value)
			c.Header("Content-Type", "application/json")
			c.JSON(200, models.UpdateMetricsModel{ID: metric.ID, MType: metric.MType, Delta: &value64})

		}
	})
}

func RegisterUpdateRoute(rg *gin.RouterGroup, repository repo.MetricsCRUDer) {

	rg.POST("/", func(c *gin.Context) {
		var metric models.UpdateMetricsModel

		var buf bytes.Buffer

		// читаем тело запроса
		_, err := buf.ReadFrom(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		}

		// десериализуем JSON
		if err = json.Unmarshal(buf.Bytes(), &metric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse	body"})
			return
		}

		if metric.MType == string(models.GaugeType) {
			err := repository.UpdateMetric(metric.ID, models.GaugeType, *metric.Value)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})
			}
		} else if metric.MType == string(models.CounterType) {
			err := repository.UpdateMetric(metric.ID, models.CounterType, *metric.Delta)
			if err == er.ErrorNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"metric name": metric.ID, "error": "metric not found"})
			}
		}

		c.JSON(200, metric)

	})

}

func RegisterMerticsRoutes(repository repo.MetricsCRUDer, logger *logger.AppLogger) *gin.Engine {

	r := gin.New()
	r.Use(RequestLogger(logger))
	r.Use((CompressionMiddleware()))
	r.Use(DeCompressionMiddleware())

	r.POST("/echo", func(c *gin.Context) {
		fmt.Println("body::", c.Request.Body)
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to read body")
			return
		}
		c.Data(http.StatusOK, "text/plain", bodyBytes)
	})

	r.LoadHTMLFiles("templates/metrics.html")

	r.GET("/", func(c *gin.Context) {
		metrics := repository.GetAllMetrics()
		fmt.Println("metrics:", metrics)
		c.HTML(http.StatusOK, "metrics.html", gin.H{
			"metrics": metrics,
		})

	})

	// TODO: rename functions!
	updateGroup := r.Group("/update")
	RegisterUpdateRoute(updateGroup, repository)
	registerUpdateGaugeRoutes(updateGroup.Group("/gauge"), repository, models.GaugeType)
	registerUpdateCounterRoutes(updateGroup.Group("/counter"), repository, models.CounterType)

	getGroup := r.Group("/value")
	RegisterGetValueRoute(getGroup, repository)
	registerGetGaugeRoutes(getGroup.Group("/gauge"), repository, models.GaugeType)
	registerGetCountRoutes(getGroup.Group("/counter"), repository, models.CounterType)

	r.GET("/ping", func(c *gin.Context) {
		data := []byte(strings.Repeat(`This is a test message`, 100))
		c.String(http.StatusOK, string(data))

		// c.String(http.StatusOK, "pong")
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
