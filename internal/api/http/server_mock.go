// Code generated by MockGen. DO NOT EDIT.
// Source: server.go

// Package http is a generated GoMock package.
package http

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockchatBotProcessor is a mock of chatBotProcessor interface.
type MockchatBotProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockchatBotProcessorMockRecorder
}

// MockchatBotProcessorMockRecorder is the mock recorder for MockchatBotProcessor.
type MockchatBotProcessorMockRecorder struct {
	mock *MockchatBotProcessor
}

// NewMockchatBotProcessor creates a new mock instance.
func NewMockchatBotProcessor(ctrl *gomock.Controller) *MockchatBotProcessor {
	mock := &MockchatBotProcessor{ctrl: ctrl}
	mock.recorder = &MockchatBotProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockchatBotProcessor) EXPECT() *MockchatBotProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockchatBotProcessor) Process(ctx context.Context, body []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", ctx, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockchatBotProcessorMockRecorder) Process(ctx, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockchatBotProcessor)(nil).Process), ctx, body)
}