// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_transaction is a generated GoMock package.
package mock_transaction

import (
	transaction "bank-api/transaction"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockITransactionService is a mock of ITransactionService interface
type MockITransactionService struct {
	ctrl     *gomock.Controller
	recorder *MockITransactionServiceMockRecorder
}

// MockITransactionServiceMockRecorder is the mock recorder for MockITransactionService
type MockITransactionServiceMockRecorder struct {
	mock *MockITransactionService
}

// NewMockITransactionService creates a new mock instance
func NewMockITransactionService(ctrl *gomock.Controller) *MockITransactionService {
	mock := &MockITransactionService{ctrl: ctrl}
	mock.recorder = &MockITransactionServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockITransactionService) EXPECT() *MockITransactionServiceMockRecorder {
	return m.recorder
}

// NewTransaction mocks base method
func (m *MockITransactionService) NewTransaction(newTransaction transaction.Transaction) (*transaction.Transaction, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewTransaction", newTransaction)
	ret0, _ := ret[0].(*transaction.Transaction)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// NewTransaction indicates an expected call of NewTransaction
func (mr *MockITransactionServiceMockRecorder) NewTransaction(newTransaction interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewTransaction", reflect.TypeOf((*MockITransactionService)(nil).NewTransaction), newTransaction)
}

// NewOperationType mocks base method
func (m *MockITransactionService) NewOperationType(newOperationType transaction.OperationType) (*transaction.OperationType, []error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewOperationType", newOperationType)
	ret0, _ := ret[0].(*transaction.OperationType)
	ret1, _ := ret[1].([]error)
	return ret0, ret1
}

// NewOperationType indicates an expected call of NewOperationType
func (mr *MockITransactionServiceMockRecorder) NewOperationType(newOperationType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewOperationType", reflect.TypeOf((*MockITransactionService)(nil).NewOperationType), newOperationType)
}
