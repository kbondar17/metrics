package routers

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"metrics/internal/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

var canGzip []string = []string{"application/json", "application/xml", "text/plain", "text/html"}

func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header

		if strings.Contains(headers.Get("Accept-Encoding"), "gzip") {
			compressWriter := CompressWriter{c.Writer}
			compressWriter.Header().Set("Content-Encoding", "gzip")
			c.Writer = compressWriter
			log.Println("sending gzip")
			c.Next()
		} else {
			log.Println("no gzip")
			c.Next()
		}
	}
}

func DeCompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header
		if !strings.Contains(headers.Get("Content-Encoding"), "gzip") {
			c.Next()
		} else {
			log.Println("decompressing gzip")

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Println("Error while reading request body: ", err)

			}

			var bu bytes.Buffer
			r, err := gzip.NewReader(bytes.NewReader(bodyBytes))
			if err != nil {
				log.Println("Error while creating gzip reader: ", err)
				return
			}

			defer r.Close()
			_, err = bu.ReadFrom(r)
			if err != nil {
				log.Println("Error while decompressing data: ", err)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bu.Bytes()))

			c.Next()
		}
	}
}
