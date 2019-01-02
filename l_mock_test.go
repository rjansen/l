package l

import (
	"github.com/stretchr/testify/mock"
)

type mockDriver struct {
	mock.Mock
}

func newMockDriver() *mockDriver {
	return new(mockDriver)
}

func (mock *mockDriver) Write(level Level, msg string, values ...Value) {
	mock.Called(level, msg, values)
}

func (mock *mockDriver) Log(level Level, msg string) LogWriter {
	args := mock.Called(level, msg)
	result := args.Get(0)
	if result == nil {
		return nil
	}
	return result.(LogWriter)
}

func (mock *mockDriver) Close() {
	mock.Called()
}

type mockLogWriter struct {
	mock.Mock
}

func newMockLogWriter() *mockLogWriter {
	return new(mockLogWriter)
}

func (mock *mockLogWriter) Write(values ...Value) {
	mock.Called(values)
}
