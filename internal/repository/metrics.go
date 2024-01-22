// функции для работы с БД
package repository

import (
	"fmt"
	"log"
	"metrics/internal/app_errors"
	models "metrics/internal/models"
)

type MetricsCRUDer interface {
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int, error)
	Create(metricName string, metricType models.MetricType) error
	GetAllMetrics() []models.MetricResponseModel
	UpdateMetric(name string, metrciType models.MetricType, value interface{}) error
}

type Storager interface {
	CheckIfMetricExists(name string, mType models.MetricType) (bool, error)
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int, error)
	Create(metricName string, metricType models.MetricType) error
	UpdateMetric(name string, metrciType models.MetricType, value interface{}) error
}

type MerticsRepo struct {
	Storage Storager
}

func NewMerticsRepo(storage Storager) MetricsCRUDer {
	return MerticsRepo{Storage: storage}
}

func (repo MerticsRepo) GetAllMetrics() []models.MetricResponseModel {
	var AllMetrics []models.MetricResponseModel

	for _, metricName := range models.MetricsDict[models.GaugeType] {
		value, err := repo.Storage.GetGaugeMetricValueByName(metricName, models.GaugeType)
		if err != nil {
			log.Printf("failed to get metric by name: %v", err)
			value = 0.0
		}

		AllMetrics = append(AllMetrics, models.MetricResponseModel{Name: metricName, Value: fmt.Sprintf("%f", value), Type: models.GaugeType})
	}

	for _, metricName := range models.MetricsDict[models.CounterType] {
		value, err := repo.Storage.GetCountMetricValueByName(metricName)
		if err != nil {
			log.Printf("failed to get metric by name: %v", err)
			value = 0
		}
		AllMetrics = append(AllMetrics, models.MetricResponseModel{Name: metricName, Value: fmt.Sprintf("%d", value), Type: models.CounterType})

	}

	return AllMetrics
}

func (repo MerticsRepo) GetCountMetricValueByName(name string) (int, error) {
	exists, err := repo.Storage.CheckIfMetricExists(name, models.CounterType)

	if !exists {
		return 0, app_errors.ErrorNotFound
	}

	if err != nil {
		log.Println("failed to get metric by name: ", err)
		return 0, err
	}
	return repo.Storage.GetCountMetricValueByName(name)
}

func (repo MerticsRepo) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	exists, err := repo.Storage.CheckIfMetricExists(name, mType)

	if !exists {
		return 0, app_errors.ErrorNotFound
	}

	if err != nil {
		log.Println("failed to get metric by name: ", err)
		return 0, err
	}
	return repo.Storage.GetGaugeMetricValueByName(name, mType)
}

func (repo MerticsRepo) Create(metricName string, metricType models.MetricType) error {
	exists, err := repo.Storage.CheckIfMetricExists(metricName, metricType)

	if err != nil {
		log.Printf("failed to check if metric exists: %v", err)
		return err
	}
	if exists {
		log.Printf("metric already exists: %v", err)
		return app_errors.AlreadyExists
	}
	log.Println("Создали метрику типа: ", metricType, " с именем: ", metricName)
	return repo.Storage.Create(metricName, metricType)

}

func (repo MerticsRepo) UpdateMetric(name string, metrciType models.MetricType, value interface{}) error {

	exists, err := repo.Storage.CheckIfMetricExists(name, metrciType)
	if err != nil {
		return err
	}

	if !exists {
		err = repo.Create(name, metrciType)
		if err != nil {
			return err
		}
	}

	return repo.Storage.UpdateMetric(name, metrciType, value)

}
