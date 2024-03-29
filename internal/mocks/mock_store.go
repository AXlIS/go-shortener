// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AXlIS/go-shortener/internal/storage (interfaces: URLWorker)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	url "github.com/AXlIS/go-shortener"
)

// MockURLWorker is a mock of URLWorker interface.
type MockURLWorker struct {
	ctrl     *gomock.Controller
	recorder *MockURLWorkerMockRecorder
}

// MockURLWorkerMockRecorder is the mock recorder for MockURLWorker.
type MockURLWorkerMockRecorder struct {
	mock *MockURLWorker
}

// NewMockURLWorker creates a new mock instance.
func NewMockURLWorker(ctrl *gomock.Controller) *MockURLWorker {
	mock := &MockURLWorker{ctrl: ctrl}
	mock.recorder = &MockURLWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLWorker) EXPECT() *MockURLWorkerMockRecorder {
	return m.recorder
}

// AddBatch mocks base method.
func (m *MockURLWorker) AddBatch(arg0 []*url.ShortenBatchInput) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddBatch", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddBatch indicates an expected call of AddBatch.
func (mr *MockURLWorkerMockRecorder) AddBatch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddBatch", reflect.TypeOf((*MockURLWorker)(nil).AddBatch), arg0)
}

// AddValue mocks base method.
func (m *MockURLWorker) AddValue(arg0, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddValue", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddValue indicates an expected call of AddValue.
func (mr *MockURLWorkerMockRecorder) AddValue(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddValue", reflect.TypeOf((*MockURLWorker)(nil).AddValue), arg0, arg1, arg2)
}

// DeleteValues mocks base method.
func (m *MockURLWorker) DeleteValues(arg0 []string, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DeleteValues", arg0, arg1)
}

// DeleteValues indicates an expected call of DeleteValues.
func (mr *MockURLWorkerMockRecorder) DeleteValues(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteValues", reflect.TypeOf((*MockURLWorker)(nil).DeleteValues), arg0, arg1)
}

// GetAllValues mocks base method.
func (m *MockURLWorker) GetAllValues(arg0 string) ([]url.Item, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllValues", arg0)
	ret0, _ := ret[0].([]url.Item)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllValues indicates an expected call of GetAllValues.
func (mr *MockURLWorkerMockRecorder) GetAllValues(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllValues", reflect.TypeOf((*MockURLWorker)(nil).GetAllValues), arg0)
}

// GetValue mocks base method.
func (m *MockURLWorker) GetValue(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValue", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValue indicates an expected call of GetValue.
func (mr *MockURLWorkerMockRecorder) GetValue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValue", reflect.TypeOf((*MockURLWorker)(nil).GetValue), arg0)
}

// Ping mocks base method.
func (m *MockURLWorker) Ping() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping.
func (mr *MockURLWorkerMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockURLWorker)(nil).Ping))
}
