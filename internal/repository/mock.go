// Code generated by MockGen. DO NOT EDIT.
// Source: metrics.go
//
// Generated by this command:
//
//	mockgen -source=metrics.go -destination=mock.go -package=repository
//

// Package repository is a generated GoMock package.
package repository

import (
	models "metrics/internal/models"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMetricsCRUDer is a mock of MetricsCRUDer interface.
type MockMetricsCRUDer struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsCRUDerMockRecorder
}

// MockMetricsCRUDerMockRecorder is the mock recorder for MockMetricsCRUDer.
type MockMetricsCRUDerMockRecorder struct {
	mock *MockMetricsCRUDer
}

// NewMockMetricsCRUDer creates a new mock instance.
func NewMockMetricsCRUDer(ctrl *gomock.Controller) *MockMetricsCRUDer {
	mock := &MockMetricsCRUDer{ctrl: ctrl}
	mock.recorder = &MockMetricsCRUDerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricsCRUDer) EXPECT() *MockMetricsCRUDerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockMetricsCRUDer) Create(metricName string, metricType models.MetricType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", metricName, metricType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockMetricsCRUDerMockRecorder) Create(metricName, metricType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockMetricsCRUDer)(nil).Create), metricName, metricType)
}

// GetAllMetrics mocks base method.
func (m *MockMetricsCRUDer) GetAllMetrics() ([]models.UpdateMetricsModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics")
	ret0, _ := ret[0].([]models.UpdateMetricsModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockMetricsCRUDerMockRecorder) GetAllMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockMetricsCRUDer)(nil).GetAllMetrics))
}

// GetCountMetricValueByName mocks base method.
func (m *MockMetricsCRUDer) GetCountMetricValueByName(name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountMetricValueByName", name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountMetricValueByName indicates an expected call of GetCountMetricValueByName.
func (mr *MockMetricsCRUDerMockRecorder) GetCountMetricValueByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountMetricValueByName", reflect.TypeOf((*MockMetricsCRUDer)(nil).GetCountMetricValueByName), name)
}

// GetGaugeMetricValueByName mocks base method.
func (m *MockMetricsCRUDer) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetricValueByName", name, mType)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetricValueByName indicates an expected call of GetGaugeMetricValueByName.
func (mr *MockMetricsCRUDerMockRecorder) GetGaugeMetricValueByName(name, mType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetricValueByName", reflect.TypeOf((*MockMetricsCRUDer)(nil).GetGaugeMetricValueByName), name, mType)
}

// Ping mocks base method.
func (m *MockMetricsCRUDer) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockMetricsCRUDerMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockMetricsCRUDer)(nil).Ping))
}

// UpdateMetric mocks base method.
func (m *MockMetricsCRUDer) UpdateMetric(name string, metrciType models.MetricType, value any, syncStorage bool, storagePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", name, metrciType, value, syncStorage, storagePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockMetricsCRUDerMockRecorder) UpdateMetric(name, metrciType, value, syncStorage, storagePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockMetricsCRUDer)(nil).UpdateMetric), name, metrciType, value, syncStorage, storagePath)
}

// UpdateMetricNew mocks base method.
func (m *MockMetricsCRUDer) UpdateMetricNew(metric models.UpdateMetricsModel, syncStorage bool, storagePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricNew", metric, syncStorage, storagePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricNew indicates an expected call of UpdateMetricNew.
func (mr *MockMetricsCRUDerMockRecorder) UpdateMetricNew(metric, syncStorage, storagePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricNew", reflect.TypeOf((*MockMetricsCRUDer)(nil).UpdateMetricNew), metric, syncStorage, storagePath)
}

// UpdateMultipleMetric mocks base method.
func (m *MockMetricsCRUDer) UpdateMultipleMetric(metrics []models.UpdateMetricsModel, syncStorage bool, storagePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMultipleMetric", metrics, syncStorage, storagePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMultipleMetric indicates an expected call of UpdateMultipleMetric.
func (mr *MockMetricsCRUDerMockRecorder) UpdateMultipleMetric(metrics, syncStorage, storagePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMultipleMetric", reflect.TypeOf((*MockMetricsCRUDer)(nil).UpdateMultipleMetric), metrics, syncStorage, storagePath)
}

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// CheckIfMetricExists mocks base method.
func (m *MockStorager) CheckIfMetricExists(name string, mType models.MetricType) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckIfMetricExists", name, mType)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckIfMetricExists indicates an expected call of CheckIfMetricExists.
func (mr *MockStoragerMockRecorder) CheckIfMetricExists(name, mType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckIfMetricExists", reflect.TypeOf((*MockStorager)(nil).CheckIfMetricExists), name, mType)
}

// Create mocks base method.
func (m *MockStorager) Create(metricName string, metricType models.MetricType) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", metricName, metricType)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockStoragerMockRecorder) Create(metricName, metricType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStorager)(nil).Create), metricName, metricType)
}

// GetAllMetrics mocks base method.
func (m *MockStorager) GetAllMetrics() ([]models.UpdateMetricsModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMetrics")
	ret0, _ := ret[0].([]models.UpdateMetricsModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMetrics indicates an expected call of GetAllMetrics.
func (mr *MockStoragerMockRecorder) GetAllMetrics() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMetrics", reflect.TypeOf((*MockStorager)(nil).GetAllMetrics))
}

// GetCountMetricValueByName mocks base method.
func (m *MockStorager) GetCountMetricValueByName(name string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCountMetricValueByName", name)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCountMetricValueByName indicates an expected call of GetCountMetricValueByName.
func (mr *MockStoragerMockRecorder) GetCountMetricValueByName(name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCountMetricValueByName", reflect.TypeOf((*MockStorager)(nil).GetCountMetricValueByName), name)
}

// GetGaugeMetricValueByName mocks base method.
func (m *MockStorager) GetGaugeMetricValueByName(name string, mType models.MetricType) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeMetricValueByName", name, mType)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeMetricValueByName indicates an expected call of GetGaugeMetricValueByName.
func (mr *MockStoragerMockRecorder) GetGaugeMetricValueByName(name, mType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeMetricValueByName", reflect.TypeOf((*MockStorager)(nil).GetGaugeMetricValueByName), name, mType)
}

// Ping mocks base method.
func (m *MockStorager) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoragerMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorager)(nil).Ping))
}

// UpdateMetric mocks base method.
func (m *MockStorager) UpdateMetric(name string, metrciType models.MetricType, value any, syncStorage bool, storagePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetric", name, metrciType, value, syncStorage, storagePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetric indicates an expected call of UpdateMetric.
func (mr *MockStoragerMockRecorder) UpdateMetric(name, metrciType, value, syncStorage, storagePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetric", reflect.TypeOf((*MockStorager)(nil).UpdateMetric), name, metrciType, value, syncStorage, storagePath)
}

// UpdateMetricNew mocks base method.
func (m *MockStorager) UpdateMetricNew(metric models.UpdateMetricsModel, syncStorage bool, storagePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMetricNew", metric, syncStorage, storagePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMetricNew indicates an expected call of UpdateMetricNew.
func (mr *MockStoragerMockRecorder) UpdateMetricNew(metric, syncStorage, storagePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMetricNew", reflect.TypeOf((*MockStorager)(nil).UpdateMetricNew), metric, syncStorage, storagePath)
}

// UpdateMultipleMetric mocks base method.
func (m *MockStorager) UpdateMultipleMetric(metrics []models.UpdateMetricsModel, syncStorage bool, storagePath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateMultipleMetric", metrics, syncStorage, storagePath)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateMultipleMetric indicates an expected call of UpdateMultipleMetric.
func (mr *MockStoragerMockRecorder) UpdateMultipleMetric(metrics, syncStorage, storagePath any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateMultipleMetric", reflect.TypeOf((*MockStorager)(nil).UpdateMultipleMetric), metrics, syncStorage, storagePath)
}
