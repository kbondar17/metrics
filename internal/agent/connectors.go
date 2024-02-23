package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	appErros "metrics/internal/errors"
	m "metrics/internal/models"

	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"go.uber.org/zap"
)

type UserClient struct {
	baseURL    string
	httpClient *gorequest.SuperAgent
	logger     *zap.SugaredLogger
}

func NewUserClient(config AgentConfig, logger *zap.SugaredLogger) UserClient {
	client := gorequest.New().
		Retry(2, 100*time.Millisecond,
			http.StatusInternalServerError,
			http.StatusBadGateway,
			http.StatusServiceUnavailable)

	client.Header.Set("Content-Type", "application/json")
	client.Client.Transport = &http.Transport{
		MaxIdleConns:    10,
		IdleConnTimeout: 30 * time.Second,
	}

	userClient := UserClient{
		httpClient: client,
		baseURL:    config.serverAddress,
		logger:     logger,
	}
	return userClient
}

func gzipCompress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("error while compressing data: %w", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("error while closing gzip writer: %w", err)
	}
	return b.Bytes(), nil
}

func (uc UserClient) SendLogsInBatches(metrics []m.UpdateMetricsModel, retrError *appErros.RetryableError) error {
	if len(metrics) == 0 {
		return fmt.Errorf("no data to send")
	}

	requestURL, err := url.JoinPath(uc.baseURL, "/updates")

	if err != nil {
		return fmt.Errorf("error while creating url:  %w", err)
	}

	bodyBytes, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("error while marshaling data:  %w", err)
	}

	compressedBody, err := gzipCompress(bodyBytes)

	if err != nil {
		return fmt.Errorf("error while compressing data:  %w", err)
	}

	_, _, errs := uc.httpClient.Post(requestURL).Set("Content-Type", "text/plain").Set("Content-Encoding", "gzip").Send(string(compressedBody)).End()
	return errors.Join(errs...)
}

func errIsRetriable(err error) bool {
	var e *url.Error
	return errors.As(err, &e)
}

func (uc UserClient) SendSingleLogCompressed(body m.UpdateMetricsModel) {
	url, err := url.JoinPath(uc.baseURL, "/update")

	if err != nil {
		uc.logger.Infoln("Error while creating url  ", err)
		return
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		uc.logger.Infow("error while marshaling data:  %w", err)
		return
	}

	compressedBody, err := gzipCompress(bodyBytes)

	if err != nil {
		uc.logger.Infow("error while compressing data:  %w", err)
		return
	}

	resp, _, errs := uc.httpClient.Post(url).Set("Content-Type", "text/plain").Set("Content-Encoding", "gzip").Send(string(compressedBody)).End()

	if errs != nil {
		uc.logger.Infoln("Error while sending data  ", errs, " response: ", resp)
		return
	}
}

func makeBody(name string, metricType m.MetricType, value string) m.UpdateMetricsModel {
	if metricType == m.GaugeType {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println("!Error while parsing float value ", err, " for metric : ", name, " value: ", value)
			value = 0
		}
		return m.UpdateMetricsModel{
			ID:    name,
			MType: string(metricType),
			Value: &value,
		}
	} else if metricType == m.CounterType {
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			log.Println("!Error while parsing int value ", err, " for metric : ", name)
			value = 0
		}
		return m.UpdateMetricsModel{
			ID:    name,
			MType: string(metricType),
			Delta: &value,
		}
	}
	return m.UpdateMetricsModel{}
}

func (uc UserClient) SendMetricContainer(data m.MetricSendContainer) {
	for metric, value := range data.GaugeMetrics {
		body := makeBody(metric, m.GaugeType, value)
		uc.SendSingleLogCompressed(body)
	}

	for metric, value := range data.UserMetrcs {
		body := makeBody(metric, m.GaugeType, value)
		uc.SendSingleLogCompressed(body)
	}

	for metric, value := range data.CounterMetrics {
		body := makeBody(metric, m.CounterType, value)
		uc.SendSingleLogCompressed(body)
	}
}

func (uc UserClient) SendMetricContainerInButches(data m.MetricSendContainer) {
	metrics := make([]m.UpdateMetricsModel, 0, len(data.GaugeMetrics)+len(data.CounterMetrics)+len(data.UserMetrcs))

	for metric, value := range data.GaugeMetrics {
		updateMetric := makeBody(metric, m.GaugeType, value)
		metrics = append(metrics, updateMetric)
	}

	for metric, value := range data.UserMetrcs {
		updateMetric := makeBody(metric, m.GaugeType, value)
		metrics = append(metrics, updateMetric)
	}
	for metric, value := range data.CounterMetrics {
		updateMetric := makeBody(metric, m.CounterType, value)
		metrics = append(metrics, updateMetric)
	}

	retrError := appErros.NewRetryableError()

	closure := func() error {
		return uc.SendLogsInBatches(metrics, retrError)
	}

	err := appErros.RetryWrapper(closure, errIsRetriable, *retrError)

	if err != nil {
		uc.logger.Infow("Error while sending data  ", err)
	}

}
