package database

import (
	"encoding/json"
	"log"
	"metrics/internal/models"
	"os"
	"strings"
	"sync"
)

func Load(fname string) ([]models.UpdateMetricsModel, error) {
	var result []models.UpdateMetricsModel
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(strings.NewReader(string(data)))

	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

var fileMutex = &sync.Mutex{}

// SaveMetric saves metric to file
func SaveMetric(fname string, data models.UpdateMetricsModel) error {
	existingData, err := Load(fname)
	if err != nil {
		log.Println("error loading file: ", fname, err)
		return err
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
		log.Println("error marshalling data: ", err)
		return err
	}
	fileMutex.Lock()
	defer fileMutex.Unlock()

	err = os.WriteFile(fname, jsonData, 0666)
	if err != nil {
		log.Println("error writing to file: ", err)
		return err
	}
	log.Println("saved metric to ", fname)
	return nil
}
