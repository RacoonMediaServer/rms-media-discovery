// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/navigator/navigator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"

	navigator "git.rms.local/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	log "github.com/apex/log"
	gomock "github.com/golang/mock/gomock"
)

// MockNavigator is a mock of Navigator interface.
type MockNavigator struct {
	ctrl     *gomock.Controller
	recorder *MockNavigatorMockRecorder
}

// MockNavigatorMockRecorder is the mock recorder for MockNavigator.
type MockNavigatorMockRecorder struct {
	mock *MockNavigator
}

// NewMockNavigator creates a new mock instance.
func NewMockNavigator(ctrl *gomock.Controller) *MockNavigator {
	mock := &MockNavigator{ctrl: ctrl}
	mock.recorder = &MockNavigatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNavigator) EXPECT() *MockNavigatorMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockNavigator) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockNavigatorMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockNavigator)(nil).Close))
}

// GetCookies mocks base method.
func (m *MockNavigator) GetCookies(urls ...string) ([]*http.Cookie, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range urls {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetCookies", varargs...)
	ret0, _ := ret[0].([]*http.Cookie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCookies indicates an expected call of GetCookies.
func (mr *MockNavigatorMockRecorder) GetCookies(urls ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCookies", reflect.TypeOf((*MockNavigator)(nil).GetCookies), urls...)
}

// NewPage mocks base method.
func (m *MockNavigator) NewPage(log *log.Entry, ctx context.Context) (navigator.Page, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewPage", log, ctx)
	ret0, _ := ret[0].(navigator.Page)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewPage indicates an expected call of NewPage.
func (mr *MockNavigatorMockRecorder) NewPage(log, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewPage", reflect.TypeOf((*MockNavigator)(nil).NewPage), log, ctx)
}
