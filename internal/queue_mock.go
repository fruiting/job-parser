// Code generated by MockGen. DO NOT EDIT.
// Source: queue.go

// Package internal is a generated GoMock package.
package internal

import (
	reflect "reflect"

	redismq "github.com/adjust/redismq"
	gomock "github.com/golang/mock/gomock"
)

// MockConsumer is a mock of Consumer interface.
type MockConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMockRecorder
}

// MockConsumerMockRecorder is the mock recorder for MockConsumer.
type MockConsumerMockRecorder struct {
	mock *MockConsumer
}

// NewMockConsumer creates a new mock instance.
func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer {
	mock := &MockConsumer{ctrl: ctrl}
	mock.recorder = &MockConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockConsumer) Get() (*redismq.Package, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get")
	ret0, _ := ret[0].(*redismq.Package)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockConsumerMockRecorder) Get() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConsumer)(nil).Get))
}

// MockTaskProcessor is a mock of TaskProcessor interface.
type MockTaskProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskProcessorMockRecorder
}

// MockTaskProcessorMockRecorder is the mock recorder for MockTaskProcessor.
type MockTaskProcessorMockRecorder struct {
	mock *MockTaskProcessor
}

// NewMockTaskProcessor creates a new mock instance.
func NewMockTaskProcessor(ctrl *gomock.Controller) *MockTaskProcessor {
	mock := &MockTaskProcessor{ctrl: ctrl}
	mock.recorder = &MockTaskProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskProcessor) EXPECT() *MockTaskProcessorMockRecorder {
	return m.recorder
}

// Ack mocks base method.
func (m *MockTaskProcessor) Ack() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ack")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ack indicates an expected call of Ack.
func (mr *MockTaskProcessorMockRecorder) Ack() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ack", reflect.TypeOf((*MockTaskProcessor)(nil).Ack))
}

// Fail mocks base method.
func (m *MockTaskProcessor) Fail() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fail")
	ret0, _ := ret[0].(error)
	return ret0
}

// Fail indicates an expected call of Fail.
func (mr *MockTaskProcessorMockRecorder) Fail() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockTaskProcessor)(nil).Fail))
}

// Requeue mocks base method.
func (m *MockTaskProcessor) Requeue() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Requeue")
	ret0, _ := ret[0].(error)
	return ret0
}

// Requeue indicates an expected call of Requeue.
func (mr *MockTaskProcessorMockRecorder) Requeue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Requeue", reflect.TypeOf((*MockTaskProcessor)(nil).Requeue))
}
