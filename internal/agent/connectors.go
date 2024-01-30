package agent

import (
	"log"
	m "metrics/internal/models"
	"net/http"
	"net/url"
	"time"

	"github.com/parnurzeal/gorequest"
)

type UserClient struct {
	baseURL    string
	httpClient *gorequest.SuperAgent
}

func NewUserClient(config AgentConfig) UserClient {

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
	}
	return userClient
}

// SendSingleLog sends single log to server
func (uc UserClient) SendSingleLog(metricName string, metricType m.MetricType, strValue string) {

	url, err := url.JoinPath(uc.baseURL, "/update/", string(metricType), metricName, strValue)

	log.Println("Sending data to:: ", url)

	if err != nil {
		log.Println("Error while creating url  ", err)
	}

	_, _, errs := uc.httpClient.Post(url).End()

	if errs != nil {
		log.Println("Error while sending data  ", errs)
	}
}

func (uc UserClient) SendMetricContainer(data m.MetricSendContainer) {

	for metric, value := range data.GaugeMetrics {
		uc.SendSingleLog(metric, m.GaugeType, value)
	}
	for metric, value := range data.UserMetrcs {
		uc.SendSingleLog(metric, m.GaugeType, value)
	}

	for metric, value := range data.CounterMetrics {
		uc.SendSingleLog(metric, m.CounterType, value)
	}

}
