package routers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	utils "metrics/internal/myutils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CompressWriter struct {
	gin.ResponseWriter
}

func (w CompressWriter) Write(b []byte) (int, error) {
	compressed, err := Compress(b)
	if err != nil {
		return 0, fmt.Errorf("failed to compress data: %v", err)
	}
	return w.ResponseWriter.Write(compressed)

}

func HashMiddleware(hashKey string, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "POST" {
			c.Next()
		}

		headers := c.Request.Header
		canonicalKey := http.CanonicalHeaderKey("Hash")

		if _, ok := headers[canonicalKey]; !ok {
			c.Next()
		}

		hash := headers.Get(canonicalKey)
		if hash == "" {
			logger.Error("Error: Hash is empty")
			c.AbortWithStatus(400)
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		// put body back so other middleware can read it
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err != nil {
			logger.Errorw("Error while reading request body", "error", err)
			c.AbortWithStatus(500)
			return
		}
		bodyHash := utils.Hash(bodyBytes, []byte(hashKey))
		if !utils.HashEqual([]byte(hash), []byte(bodyHash)) {
			logger.Info("Error: Hashes are not equal")
			c.AbortWithStatus(400)
			return
		} else {
			c.Next()
		}
	}
}

// Сведения о запросах должны содержать URI, метод запроса и время, затраченное на его выполнение.
// Сведения об ответах должны содержать код статуса и размер содержимого ответа.
func RequestLogger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var status, size int
		var start time.Time

		// generate a new request ID for each incoming request
		requestID := uuid.New().String()

		defer func() {
			if err := recover(); err != nil {
				logger.Infow("Response", "status", 500, "size", 0, "duration", time.Since(start), "requestId", requestID)
				logger.Errorw("Panic recovered", "error", err)
				c.AbortWithStatus(500)
			} else {
				logger.Infow("Response", "status", status, "size", size, "duration", time.Since(start), "requestId", requestID)
			}
		}()

		url := c.Request.URL
		method := c.Request.Method
		start = time.Now()

		logger.Infow("Request", "url", url, "method", method, "requestId", requestID)
		c.Next()
		status = c.Writer.Status()
		size = c.Writer.Size()
	}
}

// Compress compresses a slice of bytes using gzip.
func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}

	return b.Bytes(), nil
}

func CompressionMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header

		if strings.Contains(headers.Get("Accept-Encoding"), "gzip") {
			compressWriter := CompressWriter{c.Writer}
			compressWriter.Header().Set("Content-Encoding", "gzip")
			c.Writer = compressWriter
			logger.Info("sending gzip")
			c.Next()
		} else {
			logger.Info("no gzip")
			c.Next()
		}
	}
}

func DeCompressionMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		if !strings.Contains(headers.Get("Content-Encoding"), "gzip") {
			c.Next()
		} else {
			logger.Info("decompressing gzip")

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Infof("Error while reading request body: %v", err)
				return
			}

			var bu bytes.Buffer
			r, err := gzip.NewReader(bytes.NewReader(bodyBytes))
			if err != nil {
				logger.Infof("Error while creating gzip reader: %v", err)
				return
			}

			defer r.Close()
			_, err = bu.ReadFrom(r)
			if err != nil {
				logger.Infof("Error while decompressing data: %v", err)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bu.Bytes()))

			c.Next()
		}
	}
}
