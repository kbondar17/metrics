package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	appErros "metrics/internal/errors"
	m "metrics/internal/models"
	utlis "metrics/internal/myutils"
	"sync"

	"net/http"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
	"go.uber.org/zap"
)

type UserClient struct {
	baseURL    string
	httpClient *gorequest.SuperAgent
	logger     *zap.SugaredLogger
	hashKey    string
	rateLimit  int
}

func NewUserClient(config AgentConfig, logger *zap.SugaredLogger) UserClient {
	userClient := UserClient{
		baseURL:   config.serverAddress,
		logger:    logger,
		rateLimit: config.rateLimit,
	}
	return userClient
}

func newClient() *gorequest.SuperAgent {
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

	return client
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
	dataHash := ""
	if uc.hashKey != "" {
		dataHash = utlis.Hash(compressedBody, []byte(uc.hashKey))
	}

	_, _, errs := uc.httpClient.Post(requestURL).Set("Content-Type", "text/plain").Set("Hash", dataHash).Set("Content-Encoding", "gzip").Send(string(compressedBody)).End()
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
	dataHash := ""
	if uc.hashKey != "" {
		dataHash = string(utlis.Hash(compressedBody, []byte(uc.hashKey)))
	}
	client := newClient()

	resp, _, errs := client.Post(url).Set("Content-Type", "text/plain").Set("Hash", dataHash).Set("Content-Encoding", "gzip").Send(string(compressedBody)).End()

	if errs != nil {
		uc.logger.Infoln("Error while sending data  ", errs, " response: ", resp)
		return
	}

	// uc.logger.Infow("Metric sent", "ID", body.ID)
}

func (uc UserClient) SendMetricContainer(containerChan chan m.MetricSendContainer) {
	for container := range containerChan {
		metrics := container.ConvertContainerToUpdateMetricsModel()
		for _, metric := range metrics {
			uc.SendSingleLogCompressed(metric)
		}
		uc.logger.Infoln("Container sent")
	}
}

func clientWorker(jobs <-chan func()) {
	for j := range jobs {
		j()
	}
}

func (uc UserClient) SendMetricContainerWithRateLimit(dataChan chan m.MetricSendContainer) {

	jobs := make(chan func(), uc.rateLimit)

	for i := 0; i < uc.rateLimit; i++ {
		go clientWorker(jobs)
	}

	for container := range dataChan {
		startTime := time.Now()

		metrics := container.ConvertContainerToUpdateMetricsModel()
		for _, metric := range metrics {
			metric := metric
			jobs <- func() {
				uc.SendSingleLogCompressed(metric)
			}
		}

		uc.logger.Infow("Time taken to send data", "time", time.Since(startTime))
	}

}

func SplitMetricInBatches(data []m.UpdateMetricsModel, n int) [][]m.UpdateMetricsModel {
	var divided [][]m.UpdateMetricsModel

	size := len(data) / n
	for i := 0; i < n; i++ {
		start := i * size
		end := start + size
		if i == n-1 {
			end = len(data)
		}
		divided = append(divided, data[start:end])
	}

	return divided
}

func (uc UserClient) GoSend() {

}

func (uc UserClient) MetricSender(batches [][]m.UpdateMetricsModel, reportInterval int, numWorker int) {
	reportTicker := time.NewTicker(time.Duration(reportInterval) * time.Second)
	defer reportTicker.Stop()
	var wg sync.WaitGroup

	for range reportTicker.C {
		for i := 0; i < numWorker; i++ {
			wg.Add(1)
			go uc.SendMetricBatch(batches[i])
			wg.Done()
		}
	}
	wg.Wait()

}

func (uc UserClient) SendMetricBatch(metrics []m.UpdateMetricsModel) {
	uc.logger.Infow("Sending data")
	retrError := appErros.NewRetryableError()

	closure := func() error {
		return uc.SendLogsInBatches(metrics, retrError)
	}

	err := appErros.RetryWrapper(closure, errIsRetriable, *retrError)

	if err != nil {
		uc.logger.Infow("Error while sending data  ", err)
	}

}
