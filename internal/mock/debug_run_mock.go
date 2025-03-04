/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Code generated by MockGen. DO NOT EDIT.
// Source: debug_run.go
//
// Generated by this command:
//
//	mockgen -source=debug_run.go -destination=../mock/debug_run_mock.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	model "github.com/firgavin/eino-devops/internal/model"
	gomock "go.uber.org/mock/gomock"
)

// MockDebugService is a mock of DebugService interface.
type MockDebugService struct {
	ctrl     *gomock.Controller
	recorder *MockDebugServiceMockRecorder
}

// MockDebugServiceMockRecorder is the mock recorder for MockDebugService.
type MockDebugServiceMockRecorder struct {
	mock *MockDebugService
}

// NewMockDebugService creates a new mock instance.
func NewMockDebugService(ctrl *gomock.Controller) *MockDebugService {
	mock := &MockDebugService{ctrl: ctrl}
	mock.recorder = &MockDebugServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDebugService) EXPECT() *MockDebugServiceMockRecorder {
	return m.recorder
}

// CreateDebugThread mocks base method.
func (m *MockDebugService) CreateDebugThread(ctx context.Context, graphID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDebugThread", ctx, graphID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDebugThread indicates an expected call of CreateDebugThread.
func (mr *MockDebugServiceMockRecorder) CreateDebugThread(ctx, graphID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDebugThread", reflect.TypeOf((*MockDebugService)(nil).CreateDebugThread), ctx, graphID)
}

// DebugRun mocks base method.
func (m_2 *MockDebugService) DebugRun(ctx context.Context, m *model.DebugRunMeta, userInput string) (string, chan *model.NodeDebugState, chan error, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "DebugRun", ctx, m, userInput)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(chan *model.NodeDebugState)
	ret2, _ := ret[2].(chan error)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// DebugRun indicates an expected call of DebugRun.
func (mr *MockDebugServiceMockRecorder) DebugRun(ctx, m, userInput any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DebugRun", reflect.TypeOf((*MockDebugService)(nil).DebugRun), ctx, m, userInput)
}
