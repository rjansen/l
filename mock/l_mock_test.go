package mock

import (
	"context"
	"errors"
	"testing"

	"github.com/rjansen/l"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMockLogger(t *testing.T) {
	logger := NewMockLogger()
	assert.Implements(t, (*l.Logger)(nil), logger)
	logger.On("Debug", mock.Anything, mock.Anything, mock.Anything).Once()
	logger.On("Info", mock.Anything, mock.Anything, mock.Anything).Once()
	logger.On("Error", mock.Anything, mock.Anything, mock.Anything).Once()

	logger.Debug(context.Background(), "debug", l.NewValue("key", "value"))
	logger.Info(context.Background(), "info", l.NewValue("key", "value"))
	logger.Error(context.Background(), "error",
		l.NewValue("key", "value"), l.NewValue("error", errors.New("err_mock")),
	)

	logger.AssertExpectations(t)
}
