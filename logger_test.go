package l

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func clean(t assert.TestingT) {
	loggerProvider = nil
}

func setupLoggerTest(t assert.TestingT, p Provider, f FieldProvider) {
	clean(t)
	setupErr := Setup(p, f)
	assert.Nil(t, setupErr)
}

func newLoggerTest(t assert.TestingT, p Provider, f FieldProvider) Logger {
	setupLoggerTest(t, p, f)
	logger, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	return logger
}

func logTest(logger Logger) {
	logger.Debug("DebugMessage")
	logger.Debug("DebugFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
	logger.Info("InfoMessage")
	logger.Info("InfoFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
	logger.Warn("WarnMessage")
	logger.Warn("WarnFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
	logger.Error("ErrorMessage")
	logger.Error("ErrorFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
}

func TestSetupLogger(t *testing.T) {
	cases := []struct {
		provider      Provider
		fieldProvider FieldProvider
		err           error
	}{
		{newLogger, newFieldProvider(), nil},
		{nil, newFieldProvider(), ErrInvalidProvider},
		{newLogger, nil, ErrInvalidFieldProvider},
		{nil, nil, ErrInvalidProvider},
	}
	for _, c := range cases {
		clean(t)
		err := Setup(c.provider, c.fieldProvider)
		assert.Equal(t, c.err, err, fmt.Sprintf("Invalid err for data=%+v", c))
	}
}

func TestNewLogger(t *testing.T) {
	l := newLoggerTest(t, newLogger, newFieldProvider())
	logTest(l)
}

func TestOutString(t *testing.T) {
	cases := []struct {
		out Out
	}{
		{STDOUT},
		{STDERR},
		{DISCARD},
		//{OUT("invalid"), false},
	}
	for _, c := range cases {
		assert.NotNil(t, c.out.String())
	}
}

func TestOutSet(t *testing.T) {
	originalValue := "originalValue"
	var o Out
	assert.Nil(t, o.Set(originalValue))
	assert.Equal(t, originalValue, o.String())
}

func TestFormatString(t *testing.T) {
	cases := []struct {
		format Format
	}{
		{JSON},
		{TEXT},
		{Format("invalid")},
	}
	for _, c := range cases {
		assert.NotNil(t, c.format.String())
	}
}

func TestFormatSet(t *testing.T) {
	originalValue := "originalValue"
	var f Format
	assert.Nil(t, f.Set(originalValue))
	assert.Equal(t, originalValue, f.String())
}
