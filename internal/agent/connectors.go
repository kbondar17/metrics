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
	baseUrl    url.URL
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
		baseUrl:    config.ServerAddress,
		logger:     logger,
	}
	return userClient
}

// SendSingleLog sends single log to server
func (self UserClient) SendSingleLog(metricName string, metricType m.MetricType, strValue string) error {

	url, err := url.JoinPath(self.baseUrl.String(), "/update/", string(metricType), metricName, strValue)

	self.logger.Println("Sending data to:: ", url)

	if err != nil {
		err := fmt.Errorf("Error while creating url  %v", err)
		fmt.Println("Error while creating url  ", err)
		return err
	}

	_, _, errs := self.httpClient.Post(url).End()

	if errs != nil {
		err := fmt.Errorf("Error while sending data  %v", errs)
		fmt.Println("Error while sending data  ", err)
		return err
	}
	return nil
}

// SendMultipleLogsAsync iterates over map and sends logs to server
func (self UserClient) SendMultipleLogsAsync(data map[string]string, metricType m.MetricType) {
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
			if err := self.SendSingleLog(metric, metricType, value); err != nil {
				errCh <- err
			}
			<-sem // Release the token
		}(metric, value)
	}

	wg.Wait()

	self.logger.Println(metricType, " data sent")

	for err := range errCh {
		if err != nil {
			self.logger.Println("Error while sending data", err)
		}
	}

}

func (self UserClient) SendMetricContainer(data m.MetricSendContainer) {
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		defer wg.Done()
		self.SendMultipleLogsAsync(data.GaugeMetrics, m.GaugeType)
	}()

	go func() {
		defer wg.Done()
		self.SendMultipleLogsAsync(data.UserMetrcs, m.GaugeType)
	}()

	go func() {
		defer wg.Done()
		self.SendMultipleLogsAsync(data.CounterMetrics, m.CounterType)
	}()

	wg.Wait()
}
