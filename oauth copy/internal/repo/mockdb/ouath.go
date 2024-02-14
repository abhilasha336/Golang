// Code generated by MockGen. DO NOT EDIT.
// Source: oauth/internal/repo (interfaces: OauthRepoImply)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	entities "oauth/internal/entities"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockOauthRepoImply is a mock of OauthRepoImply interface.
type MockOauthRepoImply struct {
	ctrl     *gomock.Controller
	recorder *MockOauthRepoImplyMockRecorder
}

// MockOauthRepoImplyMockRecorder is the mock recorder for MockOauthRepoImply.
type MockOauthRepoImplyMockRecorder struct {
	mock *MockOauthRepoImply
}

// NewMockOauthRepoImply creates a new mock instance.
func NewMockOauthRepoImply(ctrl *gomock.Controller) *MockOauthRepoImply {
	mock := &MockOauthRepoImply{ctrl: ctrl}
	mock.recorder = &MockOauthRepoImplyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOauthRepoImply) EXPECT() *MockOauthRepoImplyMockRecorder {
	return m.recorder
}

// DeleteAndInsertRefreshToken mocks base method.
func (m *MockOauthRepoImply) DeleteAndInsertRefreshToken(arg0 context.Context, arg1, arg2, arg3, arg4 string, arg5 *string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAndInsertRefreshToken", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAndInsertRefreshToken indicates an expected call of DeleteAndInsertRefreshToken.
func (mr *MockOauthRepoImplyMockRecorder) DeleteAndInsertRefreshToken(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAndInsertRefreshToken", reflect.TypeOf((*MockOauthRepoImply)(nil).DeleteAndInsertRefreshToken), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetOauthCredentials mocks base method.
func (m *MockOauthRepoImply) GetOauthCredentials(arg0 context.Context, arg1, arg2 string) (entities.OAuthCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOauthCredentials", arg0, arg1, arg2)
	ret0, _ := ret[0].(entities.OAuthCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOauthCredentials indicates an expected call of GetOauthCredentials.
func (mr *MockOauthRepoImplyMockRecorder) GetOauthCredentials(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOauthCredentials", reflect.TypeOf((*MockOauthRepoImply)(nil).GetOauthCredentials), arg0, arg1, arg2)
}

// GetPartnerId mocks base method.
func (m *MockOauthRepoImply) GetPartnerId(arg0 context.Context, arg1, arg2 string) (string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPartnerId", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetPartnerId indicates an expected call of GetPartnerId.
func (mr *MockOauthRepoImplyMockRecorder) GetPartnerId(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPartnerId", reflect.TypeOf((*MockOauthRepoImply)(nil).GetPartnerId), arg0, arg1, arg2)
}

// GetProviderName mocks base method.
func (m *MockOauthRepoImply) GetProviderName(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProviderName", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProviderName indicates an expected call of GetProviderName.
func (mr *MockOauthRepoImplyMockRecorder) GetProviderName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProviderName", reflect.TypeOf((*MockOauthRepoImply)(nil).GetProviderName), arg0, arg1)
}

// Logout mocks base method.
func (m *MockOauthRepoImply) Logout(arg0 context.Context, arg1 entities.Refresh, arg2, arg3, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockOauthRepoImplyMockRecorder) Logout(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockOauthRepoImply)(nil).Logout), arg0, arg1, arg2, arg3, arg4)
}

// Middleware mocks base method.
func (m *MockOauthRepoImply) Middleware(arg0 context.Context, arg1 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Middleware", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Middleware indicates an expected call of Middleware.
func (mr *MockOauthRepoImplyMockRecorder) Middleware(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Middleware", reflect.TypeOf((*MockOauthRepoImply)(nil).Middleware), arg0, arg1)
}

// PostRefreshToken mocks base method.
func (m *MockOauthRepoImply) PostRefreshToken(arg0 context.Context, arg1 entities.Refresh, arg2, arg3 string, arg4 *string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostRefreshToken", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// PostRefreshToken indicates an expected call of PostRefreshToken.
func (mr *MockOauthRepoImplyMockRecorder) PostRefreshToken(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostRefreshToken", reflect.TypeOf((*MockOauthRepoImply)(nil).PostRefreshToken), arg0, arg1, arg2, arg3, arg4)
}
