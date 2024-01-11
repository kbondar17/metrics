package database

import (
	"context"
	"fmt"
	"log"
	"metrics/internal/models"
	"metrics/internal/utils"
	"strconv"

	"github.com/go-redis/redis/v8"
)

const (
	statsHashKey = "stats"
)

type RedisStorage struct {
	Client *redis.Client
	ctx    context.Context
}

// TODO: take from env
func NewRedisStorage(ctx context.Context) *RedisStorage {
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	return &RedisStorage{Client: rClient, ctx: ctx}
}

func (rdb *RedisStorage) GetCountMetricValueByName(name string) (int, error) {
	redisKey := composeRedisKey(models.CounterType, name)
	rawVal, err := rdb.Client.HGet(rdb.ctx, redisKey, "value").Result()
	if err == redis.Nil {
		return 0, utils.ErrorNotFound
	}

	value, err := strconv.ParseInt(rawVal, 10, 64)
	if err != nil {
		log.Println("failed to parse int: ", err)
		return 0, utils.ParseError
	}

	return int(value), nil

}

func (rdb *RedisStorage) GetGaugeMetricByNmae(name string) (models.GaugeMetric, error) {
	redisKey := composeRedisKey(models.GaugeType, name)
	rawVal, err := rdb.Client.HGet(rdb.ctx, redisKey, "value").Result()
	if err == redis.Nil {
		return models.GaugeMetric{}, utils.ErrorNotFound
	}

	value, err := strconv.ParseFloat(rawVal, 64)
	if err != nil {
		log.Println("failed to parse float: ", err)
		return models.GaugeMetric{}, utils.ParseError
	}

	return models.GaugeMetric{Name: name, Value: value}, nil

}

// CheckIfMetricExists checks if metric exists
func (rdb *RedisStorage) CheckIfMetricExists(name string, mType models.MetricType) (bool, error) {
	redisKey := composeRedisKey(mType, name)
	_, err := rdb.Client.HGet(rdb.ctx, redisKey, "value").Result()
	if err == redis.Nil {
		log.Println("metric not found: ", name, mType)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (rdb *RedisStorage) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	redisKey := composeRedisKey(mType, name)
	val, err := rdb.Client.HGet(rdb.ctx, redisKey, "value").Result()
	if err != nil {
		return 0, err
	}

	floatVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return floatVal, nil
}

func (rdb *RedisStorage) Create(metricName string, metricType models.MetricType) error {
	redisKey := composeRedisKey(metricType, metricName)
	err := rdb.Client.HSet(rdb.ctx, redisKey, "value", 0).Err()
	return err
}

func (rdb *RedisStorage) UpdateCounterMetric(metricToUpdate models.Metric) (int, error) {
	redisKey := composeRedisKey(metricToUpdate.Type, metricToUpdate.Name)
	err := rdb.Client.HIncrBy(rdb.ctx, redisKey, "value", int64(metricToUpdate.Value)).Err()
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (rdb *RedisStorage) UpdateMetric(name string, metrciType models.MetricType, value interface{}) error {
	redisKey := composeRedisKey(metrciType, name)
	err := rdb.Client.HSet(rdb.ctx, redisKey, "value", value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rdb *RedisStorage) UpdateGaugeMetric(name string, metrciType models.MetricType, value float64) (float64, error) {
	redisKey := composeRedisKey(metrciType, name)
	err := rdb.Client.HSet(rdb.ctx, redisKey, "value", value).Err()
	if err != nil {
		return 0, err
	}
	return value, nil
}

// composeRedisKey creates redis key
func composeRedisKey(metrciType models.MetricType, metricName string) string {
	return fmt.Sprintf("%s:%s:%s", statsHashKey, metrciType, metricName)
}
