package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"metrics/internal/models"
	"os"
	"sync"

	"go.uber.org/zap"
)

func Load(fname string, logger *zap.SugaredLogger) ([]models.UpdateMetricsModel, error) {
	var fileRMutex = &sync.RWMutex{}

	fileRMutex.RLock()
	file, err := os.Open(fname)
	if err != nil {
		fileRMutex.RUnlock()
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	defer fileRMutex.RUnlock()
	defer file.Close()

	var result []models.UpdateMetricsModel

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var metric models.UpdateMetricsModel
		err := json.Unmarshal([]byte(scanner.Text()), &metric)
		if err != nil {
			logger.Infow("error unmarshalling data: %v", err)
			continue
		}

		updated := false
		for i, m := range result {
			if m.ID == metric.ID {
				result[i] = metric
				updated = true
				break
			}
		}
		if !updated {
			result = append(result, metric)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}
	return result, nil

}

// SaveMetric saves metric to a file
func SaveMetric(fname string, metric models.UpdateMetricsModel) {
	var fileMutex = &sync.Mutex{}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Println("error opening file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(metric)
	if err != nil {
		log.Println("error marshalling data: %w", err)
	}
	_, err = file.Write(append(jsonData, '\n'))

	if err != nil {
		log.Println("error writing to file: %w", err)
	}

}
