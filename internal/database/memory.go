package database

import (
	"log"
	er "metrics/internal/errors"
	"metrics/internal/models"
	"sync"
)

type MemStorage struct {
	GaugeData map[string]float64
	CountData map[string]int64
	mu        sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{GaugeData: make(map[string]float64), CountData: make(map[string]int64)}
}

func (ms *MemStorage) CheckIfMetricExists(name string, mType models.MetricType) (bool, error) {
	switch mType {
	case models.GaugeType:
		ms.mu.RLock()
		_, ok := ms.GaugeData[name]
		ms.mu.RUnlock()
		return ok, nil
	case models.CounterType:
		ms.mu.RLock()
		_, ok := ms.CountData[name]
		ms.mu.RUnlock()
		return ok, nil
	default:
		return false, er.ErrParse
	}
}
func (ms *MemStorage) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	switch mType {
	case models.GaugeType:
		ms.mu.RLock()
		val, ok := ms.GaugeData[name]
		ms.mu.RUnlock()
		if !ok {
			return 0, er.ErrParse
		}
		return val, nil
	default:
		return 0, er.ErrParse
	}
}

func (ms *MemStorage) GetCountMetricValueByName(name string) (int64, error) {
	ms.mu.RLock()
	val, ok := ms.CountData[name]
	ms.mu.RUnlock()
	if !ok {
		return 0, er.ErrParse
	}
	return val, nil
}

func (ms *MemStorage) Create(metricName string, metricType models.MetricType) error {
	switch metricType {
	case models.GaugeType:
		ms.GaugeData[metricName] = 0
		return nil
	case models.CounterType:
		ms.CountData[metricName] = 0
		return nil
	default:
		log.Println("unknown metric type", metricType, metricName)
		return er.ErrParse
	}
}

func (ms *MemStorage) UpdateMetric(name string, metricType models.MetricType, value interface{}) error {

	log.Println("updating metric", name, metricType, value)
	switch metricType {
	case models.GaugeType:
		val, ok := value.(float64)
		if !ok {
			return er.ErrParse
		}
		ms.mu.Lock()
		ms.GaugeData[name] = val
		ms.mu.Unlock()
		return nil
	case models.CounterType:
		val, ok := value.(int64)
		if !ok {
			return er.ErrParse
		}
		ms.mu.Lock()
		ms.CountData[name] += val
		ms.mu.Unlock()
		return nil
	default:
		log.Println("Error: unknown metric type", metricType, name)
		return er.ErrParse
	}
}
