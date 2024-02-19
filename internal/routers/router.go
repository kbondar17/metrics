package routers

import (
	"io"
	"log"
	"metrics/internal/models"
	repo "metrics/internal/repository"
	get "metrics/internal/routers/get"
	post "metrics/internal/routers/post"

	_ "encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func addMiddleware(r *gin.Engine, logger *zap.SugaredLogger) {
	r.Use(RequestLogger(logger))
	r.Use(DeCompressionMiddleware(logger))
	r.Use(CompressionMiddleware(logger))
}

func RegisterMerticsRoutes(repository repo.MetricsCRUDer, logger *zap.SugaredLogger, syncStorage bool, storagePath string) *gin.Engine {
	r := gin.New()

	addMiddleware(r, logger)

	r.POST("/echo", func(c *gin.Context) {
		log.Println("body:: ", c.Request.Body)
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.String(http.StatusBadRequest, "Failed to read body")
			return
		}
		c.Data(http.StatusOK, "text/plain", bodyBytes)
	})

	r.LoadHTMLFiles("templates/metrics.html")

	r.GET("/", func(c *gin.Context) {
		metrics, err := repository.GetAllMetrics()
		log.Println("metrics:", metrics)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		c.HTML(http.StatusOK, "metrics.html", gin.H{
			"metrics": metrics,
		})

	})

	r.POST("/err", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, "text/plain")
	})

	updatesMultipleGroup := r.Group("/updates")
	post.MultipleUpdate(updatesMultipleGroup, repository, syncStorage, storagePath, logger)

	updateGroup := r.Group("/update")
	post.Update(updateGroup, repository, syncStorage, storagePath, logger)

	post.UpdateGauge(updateGroup.Group("/gauge"), repository, models.GaugeType, logger)
	post.UpdateCounter(updateGroup.Group("/counter"), repository, models.CounterType, logger)

	getGroup := r.Group("/value")
	post.GetValue(getGroup, repository, logger)

	get.GetGauge(getGroup.Group("/gauge"), repository, models.GaugeType, logger)
	get.GetCount(getGroup.Group("/counter"), repository, models.CounterType, logger)

	r.GET("/ping", func(c *gin.Context) {
		err := repository.Ping()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"failed to ping": err.Error()})
			return
		}
		c.JSON(http.StatusOK, "pong")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not Found"})
	})

	return r
}
