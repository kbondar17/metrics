package database

import (
	"encoding/json"
	"fmt"
	"log"
	"metrics/internal/models"
	"os"
	"strings"
	"sync"
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

func Load(fname string) ([]models.UpdateMetricsModel, error) {
	var fileRMutex = &sync.RWMutex{}

	fileRMutex.RLock()
	defer fileRMutex.RUnlock()

	var result []models.UpdateMetricsModel
	data, err := os.ReadFile(fname)

	if err != nil {
		log.Println("error reading file: ", err)
		if os.IsNotExist(err) {
			log.Println("file doesn't exist, creating: ", fname)
			err = createFile(fname)
			if err != nil {
				return nil, err
			}
		}
		return nil, err
	}

	if len(data) == 0 || data == nil {
		log.Println("empty file: ", fname)
		data = []byte("[]")
	}

	decoder := json.NewDecoder(strings.NewReader(string(data)))

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// SaveMetric saves metric to file
func SaveMetric(fname string, data models.UpdateMetricsModel) error {
	var fileMutex = &sync.Mutex{}

	existingData, err := Load(fname)
	if err != nil {
		return fmt.Errorf("error loading file %s: %w", fname, err)
	}
	// remove the metric to be updated
	for i, metric := range existingData {
		if metric.ID == data.ID && metric.MType == data.MType {
			existingData = append(existingData[:i], existingData[i+1:]...)
			break
		}
	}
	existingData = append(existingData, data)

	jsonData, err := json.Marshal(existingData)
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	err = os.WriteFile(fname, jsonData, 0666)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	// log.Println("saved metric to ", fname)
	return nil
}
