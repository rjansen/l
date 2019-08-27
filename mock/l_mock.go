package mock

import (
	"context"

	"github.com/rjansen/l"
	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func NewMockLogger() *MockLogger {
	return new(MockLogger)
}

func (mock *MockLogger) Debug(ctx context.Context, msg string, values ...l.Value) {
	mock.Called(ctx, msg, values)
}

func (mock *MockLogger) Info(ctx context.Context, msg string, values ...l.Value) {
	mock.Called(ctx, msg, values)
}

func (mock *MockLogger) Error(ctx context.Context, msg string, values ...l.Value) {
	mock.Called(ctx, msg, values)
}
