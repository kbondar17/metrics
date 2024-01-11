package agent

import (
	"fmt"
	"log"
	m "metrics/internal/models"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/parnurzeal/gorequest"
)

type UserClient struct {
	baseURL    url.URL
	httpClient *gorequest.SuperAgent
	logger     *log.Logger
}

var userClient UserClient

func NewUserClient(config AgentConfig, logger *log.Logger) UserClient {

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

	userClient = UserClient{
		httpClient: client,
		baseURL:    config.ServerAddress,
		logger:     logger,
	}
	return userClient
}

// SendSingleLog sends single log to server
func (uc UserClient) SendSingleLog(metricName string, metricType m.MetricType, strValue string) error {

	url, err := url.JoinPath(uc.baseURL.String(), "/update/", string(metricType), metricName, strValue)

	uc.logger.Println("Sending data to:: ", url)

	if err != nil {
		err := fmt.Errorf("error while creating url  %v", err)
		fmt.Println("Error while creating url  ", err)
		return err
	}

	_, _, errs := uc.httpClient.Post(url).End()

	if errs != nil {
		err := fmt.Errorf("error while sending data  %v", errs)
		fmt.Println("Error while sending data  ", err)
		return err
	}
	return nil
}

// SendMultipleLogsAsync iterates over map and sends logs to server
func (uc UserClient) SendMultipleLogsAsync(data map[string]string, metricType m.MetricType) {
	var wg sync.WaitGroup
	errCh := make(chan error, 50)

	//semaphore
	sem := make(chan int, 30)

	for metric, value := range data {
		wg.Add(1)
		go func(metric string, value string) {
			defer wg.Done()
			sem <- 1 // Acquire a token
			// self.SendSingleLog(metric, metricType, value)
			if err := uc.SendSingleLog(metric, metricType, value); err != nil {
				errCh <- err
			}
			<-sem // Release the token
		}(metric, value)
	}

	wg.Wait()

	uc.logger.Println(metricType, " data sent")

	for err := range errCh {
		if err != nil {
			uc.logger.Println("Error while sending data", err)
		}
	}
}

func (uc UserClient) SendMetricContainer(data m.MetricSendContainer) {
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		uc.SendMultipleLogsAsync(data.GaugeMetrics, m.GaugeType)
	}()

	go func() {
		defer wg.Done()
		uc.SendMultipleLogsAsync(data.UserMetrcs, m.GaugeType)
	}()

	go func() {
		defer wg.Done()
		uc.SendMultipleLogsAsync(data.CounterMetrics, m.CounterType)
	}()

	// wg.Wait()

}
