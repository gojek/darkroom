package service

import (
	"github.com/stretchr/testify/mock"
)

type MockManipulator struct {
	mock.Mock
}

func (m *MockManipulator) Process(spec processSpec) ([]byte, error) {
	args := m.Called(spec)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockManipulator) HasDefaultParams() bool {
	args := m.Called()
	return args.Get(0).(bool)
}
