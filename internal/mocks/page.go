// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/navigator/page.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	goquery "github.com/PuerkitoBio/goquery"
	navigator "github.com/RacoonMediaServer/rms-media-discovery/pkg/navigator"
	gomock "github.com/golang/mock/gomock"
)

// MockPage is a mock of Page interface.
type MockPage struct {
	ctrl     *gomock.Controller
	recorder *MockPageMockRecorder
}

// MockPageMockRecorder is the mock recorder for MockPage.
type MockPageMockRecorder struct {
	mock *MockPage
}

// NewMockPage creates a new mock instance.
func NewMockPage(ctrl *gomock.Controller) *MockPage {
	mock := &MockPage{ctrl: ctrl}
	mock.recorder = &MockPageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPage) EXPECT() *MockPageMockRecorder {
	return m.recorder
}

// Address mocks base method.
func (m *MockPage) Address() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Address")
	ret0, _ := ret[0].(string)
	return ret0
}

// Address indicates an expected call of Address.
func (mr *MockPageMockRecorder) Address() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Address", reflect.TypeOf((*MockPage)(nil).Address))
}

// Batch mocks base method.
func (m *MockPage) Batch(title string) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Batch", title)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// Batch indicates an expected call of Batch.
func (mr *MockPageMockRecorder) Batch(title interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Batch", reflect.TypeOf((*MockPage)(nil).Batch), title)
}

// ClearError mocks base method.
func (m *MockPage) ClearError() navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearError")
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// ClearError indicates an expected call of ClearError.
func (mr *MockPageMockRecorder) ClearError() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearError", reflect.TypeOf((*MockPage)(nil).ClearError))
}

// Click mocks base method.
func (m *MockPage) Click(selector string) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Click", selector)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// Click indicates an expected call of Click.
func (mr *MockPageMockRecorder) Click(selector interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Click", reflect.TypeOf((*MockPage)(nil).Click), selector)
}

// Close mocks base method.
func (m *MockPage) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockPageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPage)(nil).Close))
}

// Document mocks base method.
func (m *MockPage) Document() *goquery.Document {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Document")
	ret0, _ := ret[0].(*goquery.Document)
	return ret0
}

// Document indicates an expected call of Document.
func (mr *MockPageMockRecorder) Document() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Document", reflect.TypeOf((*MockPage)(nil).Document))
}

// Error mocks base method.
func (m *MockPage) Error() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(error)
	return ret0
}

// Error indicates an expected call of Error.
func (mr *MockPageMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockPage)(nil).Error))
}

// FetchContent mocks base method.
func (m *MockPage) FetchContent() navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchContent")
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// FetchContent indicates an expected call of FetchContent.
func (mr *MockPageMockRecorder) FetchContent() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchContent", reflect.TypeOf((*MockPage)(nil).FetchContent))
}

// Goto mocks base method.
func (m *MockPage) Goto(url string) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Goto", url)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// Goto indicates an expected call of Goto.
func (mr *MockPageMockRecorder) Goto(url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Goto", reflect.TypeOf((*MockPage)(nil).Goto), url)
}

// RaiseError mocks base method.
func (m *MockPage) RaiseError(err error) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RaiseError", err)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// RaiseError indicates an expected call of RaiseError.
func (mr *MockPageMockRecorder) RaiseError(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RaiseError", reflect.TypeOf((*MockPage)(nil).RaiseError), err)
}

// Sleep mocks base method.
func (m *MockPage) Sleep(d time.Duration) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Sleep", d)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// Sleep indicates an expected call of Sleep.
func (mr *MockPageMockRecorder) Sleep(d interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Sleep", reflect.TypeOf((*MockPage)(nil).Sleep), d)
}

// TracePage mocks base method.
func (m *MockPage) TracePage(fileName string) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TracePage", fileName)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// TracePage indicates an expected call of TracePage.
func (mr *MockPageMockRecorder) TracePage(fileName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TracePage", reflect.TypeOf((*MockPage)(nil).TracePage), fileName)
}

// Type mocks base method.
func (m *MockPage) Type(selector, text string) navigator.Page {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type", selector, text)
	ret0, _ := ret[0].(navigator.Page)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockPageMockRecorder) Type(selector, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockPage)(nil).Type), selector, text)
}
