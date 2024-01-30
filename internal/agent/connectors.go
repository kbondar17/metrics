package agent

import (
	"log"
	"metrics/internal/models"
	m "metrics/internal/models"
	"net/http"
	"net/url"
	"strconv"
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

// func (uc *UserClient) PostWithLogging(url string) (gorequest.Response, []error) {
// 	start := time.Now()
// 	resp, _, errs := uc.httpClient.Clone().Post(url).End()
// 	elapsed := time.Since(start)

// 	// TODO: switch to custom logger
// 	if errs != nil {
// 		log.Println("Request failed:", errs)

// 	} else {
// 		log.Printf("url: %s | time: %s | status: %s", url, elapsed, resp.Status)
// 	}
// 	return resp, errs

// }

func (uc UserClient) SendSingleLog(body models.UpdateMetricsModel) {

	url, err := url.JoinPath(uc.baseURL, "/update")

	if err != nil {
		log.Println("Error while creating url  ", err)
	}

	_, _, errs := uc.httpClient.Post(url).Set("Content-Type", "application/json").Send(body).End()

	if errs != nil {
		log.Println("Error while sending data  ", errs)
	}

}

// SendSingleLog sends single log to server
// func (uc UserClient) _SendSingleLog(metricName string, metricType m.MetricType, strValue string) {

// 	url, err := url.JoinPath(uc.baseURL, "/update/", string(metricType), metricName, strValue)

// 	// log.Println("Sending data to:: ", url)

// 	if err != nil {
// 		log.Println("Error while creating url  ", err)
// 	}

// 	// _, _, errs := uc.httpClient.Post(url).End()
// 	_, errs := uc.PostWithLogging(url)

// 	if errs != nil {
// 		log.Println("Error while sending data  ", errs)
// 	}
// }

func makeBody(name string, metricType m.MetricType, value string) (m.UpdateMetricsModel, error) {
	if metricType == m.GaugeType {
		value, err := strconv.ParseFloat(value, 64)
		return m.UpdateMetricsModel{
			ID:    name,
			MType: string(metricType),
			Value: &value,
		}, err
	} else if metricType == m.CounterType {
		value, err := strconv.ParseInt(value, 10, 64)
		return m.UpdateMetricsModel{
			ID:    name,
			MType: string(metricType),
			Delta: &value,
		}, err
	}
	return m.UpdateMetricsModel{}, nil
}

func (uc UserClient) SendMetricContainer(data m.MetricSendContainer) {

	for metric, rawValue := range data.GaugeMetrics {
		body, err := makeBody(metric, m.GaugeType, rawValue)
		if err != nil {

			log.Println("Error while parsing float64 value: ", err)

			continue
		}
		uc.SendSingleLog(body)
	}

	for metric, value := range data.UserMetrcs {
		body, err := makeBody(metric, m.GaugeType, value)
		if err != nil {

			log.Println("Error while parsing float64 value: ", err)

			continue
		}
		uc.SendSingleLog(body)
	}

	for metric, value := range data.CounterMetrics {
		body, err := makeBody(metric, m.CounterType, value)
		if err != nil {

			log.Println("Error while parsing int64 value: ", err)

			continue
		}
		uc.SendSingleLog(body)
	}

}
