package database

import (
	"context"
	"metrics/internal/models"
)

type Storager interface {
	CheckIfMetricExists(name string, mType models.MetricType) (bool, error)
	GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error)
	GetCountMetricValueByName(name string) (int, error)
	Create(metricName string, metricType models.MetricType) error
	UpdateMetric(name string, metrciType models.MetricType, value interface{}) error
}

func NewStorage(ctx context.Context) *MemStorage {
	return NewMemStorage(ctx)
	// return NewRedisStorage(ctx)
}
