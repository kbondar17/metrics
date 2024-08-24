package routers

import (
	"io"
	er "metrics/internal/errors"
	logger "metrics/internal/logger"
	"metrics/internal/models"
	"metrics/internal/repository"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestBase(t *testing.T) {
	ctrl := gomock.NewController(t)

	// mockStorage := repository.NewMockStorager(ctrl)
	// mockRepo := repository.NewMerticsRepo(mockStorage)
	mockRepo := repository.NewMockMetricsCRUDer(ctrl)

	mockRepo.EXPECT().Ping().Return(nil).AnyTimes()

	logger, err := logger.New()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	router := RegisterMerticsRoutes(mockRepo, logger, false, "/tmp/tmp.json")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "\"pong\"", w.Body.String())
}

func TestGetGaugeMetricValueByName(t *testing.T) {
	type args struct {
		url    string
		method string
		body   io.Reader
	}

	ctrl := gomock.NewController(t)

	mockRepo := repository.NewMockMetricsCRUDer(ctrl)

	mockRepo.EXPECT().GetGaugeMetricValueByName(gomock.Eq("RandomValue"), models.GaugeType).Return(12.34, nil).AnyTimes()
	mockRepo.EXPECT().GetGaugeMetricValueByName(gomock.Eq("NotExistingValue"), models.GaugeType).Return(0.0, er.ErrorNotFound).AnyTimes()
	logger, err := logger.New()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	router := RegisterMerticsRoutes(mockRepo, logger, false, "/tmp/tmp.json")

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantResponse   string
	}{
		{
			name: "Get OK",
			args: args{
				url:    "/value/gauge/RandomValue",
				method: "GET",
				body:   nil,
			},
			wantStatusCode: http.StatusOK,
			wantResponse:   "12.34",
		},
		{
			name: "Get Not Found",
			args: args{
				url:    "/value/gauge/NotExistingValue",
				method: "GET",
				body:   nil,
			},
			wantStatusCode: http.StatusNotFound,
			wantResponse:   `{"metric name":"NotExistingValue"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Equal(t, tt.wantResponse, w.Body.String())

		})
	}
}

func TestUpdateGaugeMetric(t *testing.T) {
	type args struct {
		url         string
		method      string
		body        io.Reader
		syncStorage bool
		storagePath string
	}

	ctrl := gomock.NewController(t)

	mockRepo := repository.NewMockMetricsCRUDer(ctrl)
	logger, err := logger.New()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	mockRepo.EXPECT().UpdateMetric(gomock.Eq("Alloc"), models.GaugeType, 1.1, false, "").Return(nil).AnyTimes()
	mockRepo.EXPECT().UpdateMetric(gomock.Eq("NotExistingValue"), models.GaugeType, 1.1, false, "").Return(er.ErrorNotFound).AnyTimes()

	router := RegisterMerticsRoutes(mockRepo, logger, false, "")

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantResponse   string
	}{
		{
			name: "Update OK",
			args: args{
				url:    "/update/gauge/Alloc/1.1",
				method: "POST",
				body:   nil,
			},
			wantStatusCode: http.StatusOK,
			wantResponse:   "",
		},
		{
			name: "Update Not Found",
			args: args{
				url:    "/update/gauge/NotExistingValue/1.1",
				method: "POST",
				body:   nil,
			},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   `{"metric name":"NotExistingValue"}{"error":"not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			req, _ := http.NewRequest(tt.args.method, tt.args.url, tt.args.body)
			router.ServeHTTP(w, req)
			assert.Equal(t, tt.wantStatusCode, w.Code)
			assert.Equal(t, tt.wantResponse, w.Body.String())

		})
	}
}
