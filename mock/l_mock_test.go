package mock

import (
	"errors"
	"testing"

	"github.com/rjansen/l"

	"github.com/stretchr/testify/mock"
)

func TestMockLogger(t *testing.T) {
	logger := NewMockLogger()
	logger.On("Debug", mock.Anything, mock.Anything).Once()
	logger.On("Info", mock.Anything, mock.Anything).Once()
	logger.On("Error", mock.Anything, mock.Anything).Once()
	logger.On("Close").Once()

	logger.Debug("debug", l.NewValue("key", "value"))
	logger.Info("info", l.NewValue("key", "value"))
	logger.Error("error",
		l.NewValue("key", "value"), l.NewValue("error", errors.New("err_mock")),
	)
	logger.Close()

	logger.AssertExpectations(t)
}
