package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func NewLogger(filePath string, loggerPrefix string) *log.Logger {

	// Create log dir if not exists
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multi := io.MultiWriter(file, os.Stdout)
	logger := log.New(multi, loggerPrefix, log.Ldate|log.Ltime)

	logger.Println("Logger initialized")

	return logger
}
