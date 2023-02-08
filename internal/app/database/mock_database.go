// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/blokhinnv/gophermart/internal/app/database (interfaces: Service)

// Package database is a generated GoMock package.
package database

import (
	context "context"
	reflect "reflect"

	models "github.com/blokhinnv/gophermart/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddAccrualRecord mocks base method.
func (m *MockService) AddAccrualRecord(arg0 context.Context, arg1 string, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAccrualRecord", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddAccrualRecord indicates an expected call of AddAccrualRecord.
func (mr *MockServiceMockRecorder) AddAccrualRecord(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAccrualRecord", reflect.TypeOf((*MockService)(nil).AddAccrualRecord), arg0, arg1, arg2)
}

// AddOrder mocks base method.
func (m *MockService) AddOrder(arg0 context.Context, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockServiceMockRecorder) AddOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockService)(nil).AddOrder), arg0, arg1, arg2)
}

// AddUser mocks base method.
func (m *MockService) AddUser(arg0 context.Context, arg1, arg2 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockServiceMockRecorder) AddUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockService)(nil).AddUser), arg0, arg1, arg2)
}

// FindOrderByID mocks base method.
func (m *MockService) FindOrderByID(arg0 context.Context, arg1 string) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrderByID", arg0, arg1)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrderByID indicates an expected call of FindOrderByID.
func (mr *MockServiceMockRecorder) FindOrderByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrderByID", reflect.TypeOf((*MockService)(nil).FindOrderByID), arg0, arg1)
}

// FindOrdersByUserID mocks base method.
func (m *MockService) FindOrdersByUserID(arg0 context.Context, arg1 int) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrdersByUserID", arg0, arg1)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrdersByUserID indicates an expected call of FindOrdersByUserID.
func (mr *MockServiceMockRecorder) FindOrdersByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrdersByUserID", reflect.TypeOf((*MockService)(nil).FindOrdersByUserID), arg0, arg1)
}

// FindUser mocks base method.
func (m *MockService) FindUser(arg0 context.Context, arg1, arg2 string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUser indicates an expected call of FindUser.
func (mr *MockServiceMockRecorder) FindUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUser", reflect.TypeOf((*MockService)(nil).FindUser), arg0, arg1, arg2)
}

// GetBalance mocks base method.
func (m *MockService) GetBalance(arg0 context.Context, arg1 int) (*models.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", arg0, arg1)
	ret0, _ := ret[0].(*models.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockServiceMockRecorder) GetBalance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockService)(nil).GetBalance), arg0, arg1)
}

// UpdateOrderStatus mocks base method.
func (m *MockService) UpdateOrderStatus(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrderStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrderStatus indicates an expected call of UpdateOrderStatus.
func (mr *MockServiceMockRecorder) UpdateOrderStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrderStatus", reflect.TypeOf((*MockService)(nil).UpdateOrderStatus), arg0, arg1, arg2)
}
