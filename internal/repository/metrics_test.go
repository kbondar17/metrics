package repository

import (
	"errors"
	er "metrics/internal/errors"
	"metrics/internal/logger"

	models "metrics/internal/models"
	"testing"

	"go.uber.org/mock/gomock"
)

// Общие моки для всех тестов
func setUpMockStorage(ctrl *gomock.Controller) *MockStorager {
	mockStorage := NewMockStorager(ctrl)
	mockStorage.EXPECT().CheckIfMetricExists(gomock.Eq("existing_metric"), models.GaugeType).Return(true, nil).AnyTimes()
	mockStorage.EXPECT().CheckIfMetricExists(gomock.Eq("not_existing_metric"), models.GaugeType).Return(false, nil).AnyTimes()

	mockStorage.EXPECT().GetGaugeMetricValueByName(gomock.Eq("existing_metric"), models.GaugeType).Return(2.2, nil).AnyTimes()
	mockStorage.EXPECT().GetGaugeMetricValueByName(gomock.Eq("not_existing_metric"), models.GaugeType).Return(0.0, er.ErrorNotFound).AnyTimes()

	mockStorage.EXPECT().Create(gomock.Eq("existing_metric"), models.GaugeType).Return(er.ErrAlreadyExists).AnyTimes()
	mockStorage.EXPECT().Create(gomock.Eq("not_existing_metric"), models.GaugeType).Return(nil).AnyTimes()

	return mockStorage
}

func TestMerticsRepo_GetGaugeMetricValueByName(t *testing.T) {

	type args struct {
		name  string
		mType models.MetricType
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := setUpMockStorage(ctrl)

	repo := NewMerticsRepo(mockStorage)

	tests := []struct {
		name    string
		args    args
		want    float64
		wantErr interface{}
		err     error
	}{
		{
			name:    "existing metric",
			args:    args{name: "existing_metric", mType: models.GaugeType},
			want:    2.2,
			wantErr: false,
		},

		{
			name:    "not existing metric",
			args:    args{name: "not_existing_metric", mType: models.GaugeType},
			want:    0,
			wantErr: er.ErrorNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetGaugeMetricValueByName(tt.args.name, tt.args.mType)

			if err != nil && err != tt.wantErr {
				t.Errorf("MerticsRepo.GetGaugeMetricValueByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MerticsRepo.GetGaugeMetricValueByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerticsRepo_Create(t *testing.T) {
	type args struct {
		metricName string
		metricType models.MetricType
	}

	logger, err := logger.NewAppLogger()
	if err != nil {
		t.Errorf("failed to create logger: %v", err)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockStorage := setUpMockStorage(ctrl)
	repo := NewMerticsRepo(mockStorage)

	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "create existing metric",
			args:    args{metricName: "existing_metric", metricType: models.GaugeType},
			wantErr: er.ErrAlreadyExists,
		},
		{
			name:    "create not existing metric",
			args:    args{metricName: "not_existing_metric", metricType: models.GaugeType},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.Create(tt.args.metricName, tt.args.metricType, logger); err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("MerticsRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerticsRepo_UpdateMetric(t *testing.T) {
	logger, err := logger.NewAppLogger()
	if err != nil {
		t.Errorf("failed to create logger: %v", err)
	}

	type args struct {
		name        string
		metrciType  models.MetricType
		value       interface{}
		syncStorage bool
		storagePath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "update existing metric",
			args:    args{name: "existing_metric", metrciType: models.GaugeType, value: 1, syncStorage: false, storagePath: "stub"},
			wantErr: nil,
		},

		{
			name:    "update not existing metric",
			args:    args{name: "not_existing_metric", metrciType: models.GaugeType, value: 1, syncStorage: false, storagePath: "stub"},
			wantErr: nil,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := setUpMockStorage(ctrl)
	mockStorage.EXPECT().UpdateMetric(gomock.Eq("existing_metric"), models.GaugeType, 1, false, "stub").Return(nil).AnyTimes()
	mockStorage.EXPECT().UpdateMetric(gomock.Eq("not_existing_metric"), models.GaugeType, 1, false, "stub").Return(nil).AnyTimes()

	repo := NewMerticsRepo(mockStorage)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := repo.UpdateMetric(tt.args.name, tt.args.metrciType, tt.args.value, tt.args.syncStorage, tt.args.storagePath, logger); err != nil && err != tt.wantErr {
				t.Errorf("MerticsRepo.UpdateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
