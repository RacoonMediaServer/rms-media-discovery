// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/requester/requester.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRequester is a mock of Requester interface.
type MockRequester struct {
	ctrl     *gomock.Controller
	recorder *MockRequesterMockRecorder
}

// MockRequesterMockRecorder is the mock recorder for MockRequester.
type MockRequesterMockRecorder struct {
	mock *MockRequester
}

// NewMockRequester creates a new mock instance.
func NewMockRequester(ctrl *gomock.Controller) *MockRequester {
	mock := &MockRequester{ctrl: ctrl}
	mock.recorder = &MockRequesterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequester) EXPECT() *MockRequesterMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockRequester) Download(ctx context.Context, url string) ([]byte, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, url)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Download indicates an expected call of Download.
func (mr *MockRequesterMockRecorder) Download(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockRequester)(nil).Download), ctx, url)
}

// Get mocks base method.
func (m *MockRequester) Get(ctx context.Context, url string, response interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, url, response)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockRequesterMockRecorder) Get(ctx, url, response interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRequester)(nil).Get), ctx, url, response)
}

// SetCookies mocks base method.
func (m *MockRequester) SetCookies(cookies []*http.Cookie) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCookies", cookies)
}

// SetCookies indicates an expected call of SetCookies.
func (mr *MockRequesterMockRecorder) SetCookies(cookies interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCookies", reflect.TypeOf((*MockRequester)(nil).SetCookies), cookies)
}
