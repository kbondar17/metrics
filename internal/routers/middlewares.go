package routers

import (
	"metrics/internal/logger"
	"time"

	"github.com/gin-gonic/gin"
)

// Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
// Сведения об ответах должны содержать код статуса и размер содержимого ответа.
func RequestLogger(logger *logger.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.URL
		method := c.Request.Method
		start := time.Now()
		c.Next()
		status := c.Writer.Status()
		size := c.Writer.Size()
		logger.Logger.Infow("Request: ", "url: ", url, "method: ", method, "size: ", size, "duration", time.Since(start))
		logger.Logger.Infow("Response: ", "status: ", status, "size: ", size, "duration", time.Since(start))
	}
}
