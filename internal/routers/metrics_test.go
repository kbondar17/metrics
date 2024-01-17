package routers

import (
	"io"
	"metrics/internal/database"
	"metrics/internal/models"
	"metrics/internal/repository"
	"metrics/internal/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"
	"gopkg.in/go-playground/assert.v1"
)

func TestBase(t *testing.T) {
	ctrl := gomock.NewController(t)

	mockStorage := database.NewMockStorager(ctrl)
	mockRepo := repository.NewMerticsRepo(mockStorage)
	router := RegisterMerticsRoutes(mockRepo)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
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
	mockRepo.EXPECT().GetGaugeMetricValueByName(gomock.Eq("NotExistingValue"), models.GaugeType).Return(0.0, utils.ErrorNotFound).AnyTimes()

	router := RegisterMerticsRoutes(mockRepo)

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantResponse   string
	}{
		{
			name: "Get OK",
			args: args{
				url:    "/update/gauge/RandomValue",
				method: "GET",
				body:   nil,
			},
			wantStatusCode: http.StatusOK,
			wantResponse:   "12.34",
		},
		{
			name: "Get Not Found",
			args: args{
				url:    "/update/gauge/NotExistingValue",
				method: "GET",
				body:   nil,
			},
			wantStatusCode: http.StatusBadRequest,
			wantResponse:   `{"error":"metric not found","metric name":"NotExistingValue"}`,
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
		url    string
		method string
		body   io.Reader
	}

	ctrl := gomock.NewController(t)

	mockRepo := repository.NewMockMetricsCRUDer(ctrl)

	mockRepo.EXPECT().UpdateMetric(gomock.Eq("Alloc"), models.GaugeType, 1.1).Return(nil).AnyTimes()
	mockRepo.EXPECT().UpdateMetric(gomock.Eq("NotExistingValue"), models.GaugeType, 1.1).Return(utils.ErrorNotFound).AnyTimes()

	router := RegisterMerticsRoutes(mockRepo)

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
			wantResponse:   `{"error":"metric not found","metric name":"NotExistingValue"}`,
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
