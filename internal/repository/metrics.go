// функции для работы с БД
package repository

import (
	"fmt"
	er "metrics/internal/errors"

	models "metrics/internal/models"

	"go.uber.org/zap"
)

type MetricsCRUDer interface {
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int64, error)
	Create(metricName string, metricType models.MetricType) error
	GetAllMetrics() ([]models.UpdateMetricsModel, error)
	UpdateMetric(name string, metrciType models.MetricType, value interface{}, syncStorage bool, storagePath string) error
	UpdateMultipleMetric(metrics []models.UpdateMetricsModel) error
	Ping() error
}

type Storager interface {
	CheckIfMetricExists(name string, mType models.MetricType) (bool, error)
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int64, error)
	Create(metricName string, metricType models.MetricType) error
	UpdateMetric(name string, metrciType models.MetricType, value interface{}, syncStorage bool, storagePath string) error
	GetAllMetrics() ([]models.UpdateMetricsModel, error)
	UpdateMultipleMetric(metrics []models.UpdateMetricsModel) error
	Ping() error
}

type MerticsRepo struct {
	Storage Storager
	logger  *zap.SugaredLogger
}

func NewMerticsRepo(storage Storager, logger *zap.SugaredLogger) MetricsCRUDer {
	return MerticsRepo{Storage: storage, logger: logger}
}

func (repo MerticsRepo) Ping() error {
	err := repo.Storage.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping storage %w", err)
	}
	return nil
}

func (repo MerticsRepo) UpdateMultipleMetric(metrics []models.UpdateMetricsModel) error {
	return repo.Storage.UpdateMultipleMetric(metrics)
}

func (repo MerticsRepo) GetAllMetrics() ([]models.UpdateMetricsModel, error) {
	metrics, err := repo.Storage.GetAllMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to get all metrics: %w", err)
	}

	return metrics, err
}

func (repo MerticsRepo) GetCountMetricValueByName(name string) (int64, error) {
	exists, err := repo.Storage.CheckIfMetricExists(name, models.CounterType)

	if !exists {
		return 0, er.ErrorNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get metric by name: %w", err)
	}
	return repo.Storage.GetCountMetricValueByName(name)
}

func (repo MerticsRepo) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	exists, err := repo.Storage.CheckIfMetricExists(name, mType)

	if !exists {
		return 0, er.ErrorNotFound
	}

	if err != nil {
		return 0, fmt.Errorf("failed to get metric by name: %w", err)
	}
	return repo.Storage.GetGaugeMetricValueByName(name, mType)
}

func (repo MerticsRepo) Create(metricName string, metricType models.MetricType) error {
	exists, err := repo.Storage.CheckIfMetricExists(metricName, metricType)

	if err != nil {
		return fmt.Errorf("failed to check if metric exists: %w", err)
	}
	if exists {
		return er.ErrAlreadyExists
	}
	repo.logger.Infof("Создали метрику типа: ", metricType, " с именем: ", metricName)
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
