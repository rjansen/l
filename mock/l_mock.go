package mock

import (
	"github.com/rjansen/l"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func NewMockLogger() *MockLogger {
	return new(MockLogger)
}

func (mock *MockLogger) Debug(msg string, values ...l.Value) {
	mock.Called(msg, values)
}

func (mock *MockLogger) Info(msg string, values ...l.Value) {
	mock.Called(msg, values)
}

func (mock *MockLogger) Error(msg string, values ...l.Value) {
	mock.Called(msg, values)
}

func (mock *MockLogger) Close() {
	mock.Called()
}
