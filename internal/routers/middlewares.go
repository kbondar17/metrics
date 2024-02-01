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
	fmt.Println("data size before compression: ", len(b))
	compressed, err := Compress(b)
	if err != nil {
		return 0, fmt.Errorf("failed to compress data: %v", err)
	}
	fmt.Println("data size after compression: ", len(compressed))
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

// Decompress распаковывает слайс байт.
func Decompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		fmt.Println("data is empty")
		return nil, fmt.Errorf("data is empty")
	}

	var b bytes.Buffer
	r, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("failed init decompress reader: %v", err)
	}

	defer r.Close()
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, fmt.Errorf("failed decompress data: %v", err)
	}
	fmt.Println("data:: ", data)
	return b.Bytes(), nil
}

var canGzip []string = []string{"application/json", "application/xml", "text/plain", "text/html"}

func haveCommonElement(a, b []string) bool {
	for _, v := range a {
		for _, w := range b {
			fmt.Println("comparing: ", "conType: ", v, "cat zip: ", w)
			if v == w {
				return true
			}
		}
	}
	return false
}

func CompressionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		headers := c.Request.Header

		// contTypes := strings.Split(headers.Get("Accept"), ",")
		// && haveCommonElement(contTypes, canGzip)
		if strings.Contains(headers.Get("Accept-Encoding"), "gzip") {
			compressWriter := CompressWriter{c.Writer}
			compressWriter.Header().Set("Content-Encoding", "gzip")
			c.Writer = compressWriter
			log.Println("sending gzip")
			c.Next()
			return
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
			log.Println("обрабатываем gzip")

			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				fmt.Println("Error while reading request body: ", err)

			}
			// fmt.Println("bodyBytes:", bodyBytes)
			// fmt.Println("string(bodyBytes):", string(bodyBytes))

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

			fmt.Println("decompressed data:: ", bu.String())

			c.Request.Body = io.NopCloser(bytes.NewBuffer(bu.Bytes()))

			c.Next()
		}
	}
}
