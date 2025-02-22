// Code generated by MockGen. DO NOT EDIT.
// Source: internal/adapters/repository/company.go
//
// Generated by this command:
//
//	mockgen -source=internal/adapters/repository/company.go -package=mock_mysql -destination=test/mock/repository/company.go
//

// Package mock_mysql is a generated GoMock package.
package mock_mysql

import (
	context "context"
	reflect "reflect"

	domain "github.com/rhuandantas/xm-challenge/internal/core/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockCompanyRepo is a mock of CompanyRepo interface.
type MockCompanyRepo struct {
	ctrl     *gomock.Controller
	recorder *MockCompanyRepoMockRecorder
	isgomock struct{}
}

// MockCompanyRepoMockRecorder is the mock recorder for MockCompanyRepo.
type MockCompanyRepoMockRecorder struct {
	mock *MockCompanyRepo
}

// NewMockCompanyRepo creates a new mock instance.
func NewMockCompanyRepo(ctrl *gomock.Controller) *MockCompanyRepo {
	mock := &MockCompanyRepo{ctrl: ctrl}
	mock.recorder = &MockCompanyRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompanyRepo) EXPECT() *MockCompanyRepoMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCompanyRepo) Create(ctx context.Context, company *domain.Company) (*domain.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, company)
	ret0, _ := ret[0].(*domain.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCompanyRepoMockRecorder) Create(ctx, company any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCompanyRepo)(nil).Create), ctx, company)
}

// DeleteByID mocks base method.
func (m *MockCompanyRepo) DeleteByID(ctx context.Context, id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID.
func (mr *MockCompanyRepoMockRecorder) DeleteByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockCompanyRepo)(nil).DeleteByID), ctx, id)
}

// GetByID mocks base method.
func (m *MockCompanyRepo) GetByID(ctx context.Context, id string) (*domain.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(*domain.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCompanyRepoMockRecorder) GetByID(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCompanyRepo)(nil).GetByID), ctx, id)
}

// GetByName mocks base method.
func (m *MockCompanyRepo) GetByName(ctx context.Context, name string) (*domain.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", ctx, name)
	ret0, _ := ret[0].(*domain.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockCompanyRepoMockRecorder) GetByName(ctx, name any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockCompanyRepo)(nil).GetByName), ctx, name)
}

// Update mocks base method.
func (m *MockCompanyRepo) Update(ctx context.Context, id string, company *domain.Company) (*domain.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, company)
	ret0, _ := ret[0].(*domain.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCompanyRepoMockRecorder) Update(ctx, id, company any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCompanyRepo)(nil).Update), ctx, id, company)
}
