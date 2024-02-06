package agent

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"log"
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

func (uc UserClient) SendSingleLogCompressed(body m.UpdateMetricsModel) {

	url, err := url.JoinPath(uc.baseURL, "/update")

	if err != nil {
		log.Println("Error while creating url  ", err)
		return
	}

	bodyBytes, err := json.Marshal(body)

	if err != nil {
		log.Println("Error while marshaling body  ", err)
		return
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)

	if _, err := gz.Write(bodyBytes); err != nil {
		log.Println("Error while writing to gzip writer: ", err)
		return
	}

	if err := gz.Close(); err != nil {
		log.Println("Error while closing gzip writer: ", err)
		return
	}

	compressedBody := b.Bytes()
	resp, _, errs := uc.httpClient.Post(url).Set("Content-Type", "text/plain").Set("Content-Encoding", "gzip").Send(string(compressedBody)).End()

	if errs != nil {
		log.Println("Error while sending data  ", errs, " response: ", resp)
		return
	}
}

func (uc UserClient) SendSingleLog(body m.UpdateMetricsModel) {

	url, err := url.JoinPath(uc.baseURL, "/update")

	if err != nil {
		log.Println("Error while creating url  ", err)
	}

	_, _, errs := uc.httpClient.Post(url).Set("Content-Type", "application/json").Send(body).End()

	if errs != nil {
		log.Println("Error while sending data  ", errs)
	}

}

func makeBody(name string, metricType m.MetricType, value string) m.UpdateMetricsModel {
	if metricType == m.GaugeType {
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			log.Println("!Error while parsing float value ", err, " for metric : ", name)
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
		uc.SendSingleLog(body)
	}

	for metric, value := range data.CounterMetrics {
		body := makeBody(metric, m.CounterType, value)
		uc.SendSingleLogCompressed(body)
	}

}
