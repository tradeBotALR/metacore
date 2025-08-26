package mocks

import (
	"reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRowInterface is a mock of RowInterface interface.
type MockRowInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRowInterfaceMockRecorder
}

// MockRowInterfaceMockRecorder is the mock recorder for MockRowInterface.
type MockRowInterfaceMockRecorder struct {
	mock *MockRowInterface
}

// NewMockRowInterface creates a new mock instance.
func NewMockRowInterface(ctrl *gomock.Controller) *MockRowInterface {
	mock := &MockRowInterface{ctrl: ctrl}
	mock.recorder = &MockRowInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRowInterface) EXPECT() *MockRowInterfaceMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockRowInterface) Scan(dest ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range dest {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockRowInterfaceMockRecorder) Scan(dest ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range dest {
		varargs = append(varargs, a)
	}
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockRowInterface)(nil).Scan), varargs...)
}
