package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"metrics/internal/models"
	"os"
	"sync"

	"go.uber.org/zap"
)

func createFile(fname string) error {
	file, err := os.OpenFile(fname, os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte("[]"))
	if err != nil {
		return err
	}
	return nil
}

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

// TODO: LOGGER
// SaveMetric saves metric to file
func SaveMetric(fname string, data models.UpdateMetricsModel) error {
	var fileMutex = &sync.Mutex{}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	file, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}
	_, err = file.Write(append(jsonData, '\n'))

	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
