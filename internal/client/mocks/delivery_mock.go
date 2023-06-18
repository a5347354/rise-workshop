// Code generated by MockGen. DO NOT EDIT.
// Source: ./delivery.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockMetrics is a mock of Metrics interface.
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics.
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance.
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// FailTotal mocks base method.
func (m *MockMetrics) FailTotal(lvs ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range lvs {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "FailTotal", varargs...)
}

// FailTotal indicates an expected call of FailTotal.
func (mr *MockMetricsMockRecorder) FailTotal(lvs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FailTotal", reflect.TypeOf((*MockMetrics)(nil).FailTotal), lvs...)
}

// ProcessDuration mocks base method.
func (m *MockMetrics) ProcessDuration(t time.Time, lvs ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{t}
	for _, a := range lvs {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "ProcessDuration", varargs...)
}

// ProcessDuration indicates an expected call of ProcessDuration.
func (mr *MockMetricsMockRecorder) ProcessDuration(t interface{}, lvs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{t}, lvs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessDuration", reflect.TypeOf((*MockMetrics)(nil).ProcessDuration), varargs...)
}

// SuccessTotal mocks base method.
func (m *MockMetrics) SuccessTotal(lvs ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range lvs {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SuccessTotal", varargs...)
}

// SuccessTotal indicates an expected call of SuccessTotal.
func (mr *MockMetricsMockRecorder) SuccessTotal(lvs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SuccessTotal", reflect.TypeOf((*MockMetrics)(nil).SuccessTotal), lvs...)
}

// WebsocketConnectionNumber mocks base method.
func (m *MockMetrics) WebsocketConnectionNumber(arg0 ...string) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "WebsocketConnectionNumber", varargs...)
}

// WebsocketConnectionNumber indicates an expected call of WebsocketConnectionNumber.
func (mr *MockMetricsMockRecorder) WebsocketConnectionNumber(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WebsocketConnectionNumber", reflect.TypeOf((*MockMetrics)(nil).WebsocketConnectionNumber), arg0...)
}