package database

import (
	"encoding/json"
	"metrics/internal/models"
	"os"
	"strings"
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

func Save(fname string, data []models.UpdateMetricsModel) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(fname, jsonData, 0666)
}
