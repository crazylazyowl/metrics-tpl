// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics (interfaces: MetricRegistry,MetricFetcher,MetricUpdater)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	metrics "github.com/crazylazyowl/metrics-tpl/internal/usecase/metrics"
	gomock "github.com/golang/mock/gomock"
)

// MockMetricRegistry is a mock of MetricRegistry interface.
type MockMetricRegistry struct {
	ctrl     *gomock.Controller
	recorder *MockMetricRegistryMockRecorder
}

// MockMetricRegistryMockRecorder is the mock recorder for MockMetricRegistry.
type MockMetricRegistryMockRecorder struct {
	mock *MockMetricRegistry
}

// NewMockMetricRegistry creates a new mock instance.
func NewMockMetricRegistry(ctrl *gomock.Controller) *MockMetricRegistry {
	mock := &MockMetricRegistry{ctrl: ctrl}
	mock.recorder = &MockMetricRegistryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricRegistry) EXPECT() *MockMetricRegistryMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockMetricRegistry) Fetch(arg0 context.Context) ([]metrics.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0)
	ret0, _ := ret[0].([]metrics.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockMetricRegistryMockRecorder) Fetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockMetricRegistry)(nil).Fetch), arg0)
}

// FetchOne mocks base method.
func (m *MockMetricRegistry) FetchOne(arg0 context.Context, arg1 metrics.Metric) (metrics.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOne", arg0, arg1)
	ret0, _ := ret[0].(metrics.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOne indicates an expected call of FetchOne.
func (mr *MockMetricRegistryMockRecorder) FetchOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOne", reflect.TypeOf((*MockMetricRegistry)(nil).FetchOne), arg0, arg1)
}

// Update mocks base method.
func (m *MockMetricRegistry) Update(arg0 context.Context, arg1 []metrics.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockMetricRegistryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMetricRegistry)(nil).Update), arg0, arg1)
}

// UpdateOne mocks base method.
func (m *MockMetricRegistry) UpdateOne(arg0 context.Context, arg1 metrics.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOne", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOne indicates an expected call of UpdateOne.
func (mr *MockMetricRegistryMockRecorder) UpdateOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOne", reflect.TypeOf((*MockMetricRegistry)(nil).UpdateOne), arg0, arg1)
}

// MockMetricFetcher is a mock of MetricFetcher interface.
type MockMetricFetcher struct {
	ctrl     *gomock.Controller
	recorder *MockMetricFetcherMockRecorder
}

// MockMetricFetcherMockRecorder is the mock recorder for MockMetricFetcher.
type MockMetricFetcherMockRecorder struct {
	mock *MockMetricFetcher
}

// NewMockMetricFetcher creates a new mock instance.
func NewMockMetricFetcher(ctrl *gomock.Controller) *MockMetricFetcher {
	mock := &MockMetricFetcher{ctrl: ctrl}
	mock.recorder = &MockMetricFetcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricFetcher) EXPECT() *MockMetricFetcherMockRecorder {
	return m.recorder
}

// Fetch mocks base method.
func (m *MockMetricFetcher) Fetch(arg0 context.Context) ([]metrics.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0)
	ret0, _ := ret[0].([]metrics.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockMetricFetcherMockRecorder) Fetch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockMetricFetcher)(nil).Fetch), arg0)
}

// FetchOne mocks base method.
func (m *MockMetricFetcher) FetchOne(arg0 context.Context, arg1 metrics.Metric) (metrics.Metric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchOne", arg0, arg1)
	ret0, _ := ret[0].(metrics.Metric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchOne indicates an expected call of FetchOne.
func (mr *MockMetricFetcherMockRecorder) FetchOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchOne", reflect.TypeOf((*MockMetricFetcher)(nil).FetchOne), arg0, arg1)
}

// MockMetricUpdater is a mock of MetricUpdater interface.
type MockMetricUpdater struct {
	ctrl     *gomock.Controller
	recorder *MockMetricUpdaterMockRecorder
}

// MockMetricUpdaterMockRecorder is the mock recorder for MockMetricUpdater.
type MockMetricUpdaterMockRecorder struct {
	mock *MockMetricUpdater
}

// NewMockMetricUpdater creates a new mock instance.
func NewMockMetricUpdater(ctrl *gomock.Controller) *MockMetricUpdater {
	mock := &MockMetricUpdater{ctrl: ctrl}
	mock.recorder = &MockMetricUpdaterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetricUpdater) EXPECT() *MockMetricUpdaterMockRecorder {
	return m.recorder
}

// Update mocks base method.
func (m *MockMetricUpdater) Update(arg0 context.Context, arg1 []metrics.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockMetricUpdaterMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMetricUpdater)(nil).Update), arg0, arg1)
}

// UpdateOne mocks base method.
func (m *MockMetricUpdater) UpdateOne(arg0 context.Context, arg1 metrics.Metric) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOne", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOne indicates an expected call of UpdateOne.
func (mr *MockMetricUpdaterMockRecorder) UpdateOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOne", reflect.TypeOf((*MockMetricUpdater)(nil).UpdateOne), arg0, arg1)
}
