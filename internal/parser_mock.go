// Code generated by MockGen. DO NOT EDIT.
// Source: parser.go

// Package internal is a generated GoMock package.
package internal

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJobsParser is a mock of JobsParser interface.
type MockJobsParser struct {
	ctrl     *gomock.Controller
	recorder *MockJobsParserMockRecorder
}

// MockJobsParserMockRecorder is the mock recorder for MockJobsParser.
type MockJobsParserMockRecorder struct {
	mock *MockJobsParser
}

// NewMockJobsParser creates a new mock instance.
func NewMockJobsParser(ctrl *gomock.Controller) *MockJobsParser {
	mock := &MockJobsParser{ctrl: ctrl}
	mock.recorder = &MockJobsParserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJobsParser) EXPECT() *MockJobsParserMockRecorder {
	return m.recorder
}

// GeneralLink mocks base method.
func (m *MockJobsParser) GeneralLink(position Name) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneralLink", position)
	ret0, _ := ret[0].(string)
	return ret0
}

// GeneralLink indicates an expected call of GeneralLink.
func (mr *MockJobsParserMockRecorder) GeneralLink(position interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneralLink", reflect.TypeOf((*MockJobsParser)(nil).GeneralLink), position)
}

// ItemsCount mocks base method.
func (m *MockJobsParser) ItemsCount() uint16 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ItemsCount")
	ret0, _ := ret[0].(uint16)
	return ret0
}

// ItemsCount indicates an expected call of ItemsCount.
func (mr *MockJobsParserMockRecorder) ItemsCount() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ItemsCount", reflect.TypeOf((*MockJobsParser)(nil).ItemsCount))
}

// Links mocks base method.
func (m *MockJobsParser) Links(dom string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Links", dom)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Links indicates an expected call of Links.
func (mr *MockJobsParserMockRecorder) Links(dom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Links", reflect.TypeOf((*MockJobsParser)(nil).Links), dom)
}

// PagesCount mocks base method.
func (m *MockJobsParser) PagesCount(dom string) (uint16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PagesCount", dom)
	ret0, _ := ret[0].(uint16)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PagesCount indicates an expected call of PagesCount.
func (mr *MockJobsParserMockRecorder) PagesCount(dom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PagesCount", reflect.TypeOf((*MockJobsParser)(nil).PagesCount), dom)
}

// ParseDetail mocks base method.
func (m *MockJobsParser) ParseDetail(dom string) (*Job, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseDetail", dom)
	ret0, _ := ret[0].(*Job)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseDetail indicates an expected call of ParseDetail.
func (mr *MockJobsParserMockRecorder) ParseDetail(dom interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseDetail", reflect.TypeOf((*MockJobsParser)(nil).ParseDetail), dom)
}

// Parser mocks base method.
func (m *MockJobsParser) Parser() Parser {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Parser")
	ret0, _ := ret[0].(Parser)
	return ret0
}

// Parser indicates an expected call of Parser.
func (mr *MockJobsParserMockRecorder) Parser() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Parser", reflect.TypeOf((*MockJobsParser)(nil).Parser))
}

// SearchPageLink mocks base method.
func (m *MockJobsParser) SearchPageLink(pageNumber uint16) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchPageLink", pageNumber)
	ret0, _ := ret[0].(string)
	return ret0
}

// SearchPageLink indicates an expected call of SearchPageLink.
func (mr *MockJobsParserMockRecorder) SearchPageLink(pageNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchPageLink", reflect.TypeOf((*MockJobsParser)(nil).SearchPageLink), pageNumber)
}
