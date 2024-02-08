// функции для работы с БД
package repository

import (
	"log"
	er "metrics/internal/errors"

	models "metrics/internal/models"
)

type MetricsCRUDer interface {
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int64, error)
	Create(metricName string, metricType models.MetricType) error
	GetAllMetrics() []models.UpdateMetricsModel
	UpdateMetric(name string, metrciType models.MetricType, value interface{}, syncStorage bool, storagePath string) error
}

type Storager interface {
	CheckIfMetricExists(name string, mType models.MetricType) (bool, error)
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int64, error)
	Create(metricName string, metricType models.MetricType) error
	UpdateMetric(name string, metrciType models.MetricType, value interface{}, syncStorage bool, storagePath string) error
	GetAllMetrics() []models.UpdateMetricsModel
}

type MerticsRepo struct {
	Storage Storager
}

func NewMerticsRepo(storage Storager) MetricsCRUDer {
	return MerticsRepo{Storage: storage}
}

func (repo MerticsRepo) GetAllMetrics() []models.UpdateMetricsModel {
	metrics := repo.Storage.GetAllMetrics()
	log.Println("all metrics: ", metrics)
	return metrics
}

func (repo MerticsRepo) GetCountMetricValueByName(name string) (int64, error) {
	exists, err := repo.Storage.CheckIfMetricExists(name, models.CounterType)

	if !exists {
		return 0, er.ErrorNotFound
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
		return 0, er.ErrorNotFound
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
		// log.Printf("metric already exists: %v", err)
		return er.ErrAlreadyExists
	}
	log.Println("Создали метрику типа: ", metricType, " с именем: ", metricName)
	return repo.Storage.Create(metricName, metricType)

}

func (repo MerticsRepo) UpdateMetric(name string, metrciType models.MetricType, value interface{}, syncStorage bool, storagePath string) error {
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
	return repo.Storage.UpdateMetric(name, metrciType, value, syncStorage, storagePath)
}
