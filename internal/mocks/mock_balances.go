package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/Elvilius/go-musthave-diploma-tpl.git/internal/models"
	gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockBalancesStore struct {
	ctrl     *gomock.Controller
	recorder *MockBalancesStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockBalancesStoreMockRecorder struct {
	mock *MockBalancesStore
}

// NewMockStore creates a new mock instance.
func NewMockBalancesStore(ctrl *gomock.Controller) *MockBalancesStore {
	mock := &MockBalancesStore{ctrl: ctrl}
	mock.recorder = &MockBalancesStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBalancesStore) EXPECT() *MockBalancesStoreMockRecorder {
	return m.recorder
}

// GetBalance mocks base method.
func (m *MockBalancesStore) GetBalance(ctx context.Context, userID uint64) (models.Balance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBalance", ctx, userID)
	ret0, _ := ret[0].(models.Balance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBalance indicates an expected call of GetBalance.
func (mr *MockBalancesStoreMockRecorder) GetBalance(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBalance", reflect.TypeOf((*MockBalancesStore)(nil).GetBalance), ctx, userID)
}

// GetWithdraws mocks base method.
func (m *MockBalancesStore) GetWithdraws(ctx context.Context, userID uint64) ([]models.GetWithdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdraws", ctx, userID)
	ret0, _ := ret[0].([]models.GetWithdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdraws indicates an expected call of GetWithdraws.
func (mr *MockBalancesStoreMockRecorder) GetWithdraws(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdraws", reflect.TypeOf((*MockBalancesStore)(nil).GetWithdraws), ctx, userID)
}

// Withdraw mocks base method.
func (m *MockBalancesStore) Withdraw(ctx context.Context, userID uint64, order string, sum float32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", ctx, userID, order, sum)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockBalancesStoreMockRecorder) Withdraw(ctx, userID, order, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockBalancesStore)(nil).Withdraw), ctx, userID, order, sum)
}
