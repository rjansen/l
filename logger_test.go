package logger

import (
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
)

var (
	configTest = Configuration{Provider: LOGRUS, Format: TEXTColor}
)

func setupLoggerTest(t assert.TestingT, config Configuration) {
	setupErr := Setup(config)
	assert.Nil(t, setupErr)
}

func newLoggerTest(t assert.TestingT, config Configuration) Logger {
	setupLoggerTest(t, config)
	logger := NewLogger()
	assert.NotNil(t, logger)
	return logger
}

func newLogrusByConfigTest(t assert.TestingT, config Configuration) Logger {
	logger := NewLoggerByConfig(config)
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
		config  Configuration
		success bool
	}{
		{Configuration{Out: STDOUT}, false},
		{Configuration{Out: STDERR}, false},
		{Configuration{Out: DISCARD}, false},
		{Configuration{Provider: LOGRUS, Out: STDOUT}, true},
		{Configuration{Provider: ZAP, Out: STDOUT}, true},
		{Configuration{Provider: LOGRUS, Out: STDERR}, true},
		{Configuration{Provider: ZAP, Out: STDERR}, true},
		{Configuration{Provider: LOGRUS, Out: DISCARD}, true},
		{Configuration{Provider: ZAP, Out: DISCARD}, true},
		{Configuration{Provider: LOGRUS, Out: Out("/tmp/glive_test.log")}, true},
		{Configuration{Provider: ZAP, Out: Out("/tmp/glive_test.log")}, true},
		//{Configuration{Provider: LOGRUS, Out: Out("%$$/{}")}, false},
		//{Configuration{Provider: ZAP, Out: Out("$$%/{}")}, false},
		{Configuration{Provider: LOGRUS, Format: JSON}, true},
		{Configuration{Provider: ZAP, Format: JSON}, true},
		{Configuration{Provider: LOGRUS, Format: TEXT}, true},
		{Configuration{Provider: ZAP, Format: TEXT}, true},
		{Configuration{Provider: LOGRUS, Format: TEXTColor}, true},
		{Configuration{Provider: ZAP, Format: TEXTColor}, true},
		{Configuration{Provider: LOGRUS, Format: JSONColor}, true},
		{Configuration{Provider: ZAP, Format: JSONColor}, true},
		{Configuration{Provider: LOGRUS, Level: Level(13)}, true},
		{Configuration{Provider: ZAP, Level: Level(13)}, true},
		//{Configuration{Provider: LOGRUS, Level: Level(99)}, false},
		//{Configuration{Provider: ZAP, Level: Level(99)}, false},
	}
	for _, c := range cases {
		err := Setup(c.config)
		if c.success {
			assert.Nil(t, err, fmt.Sprintf("Not Nil for data=%+v", c))
		} else {
			assert.NotNil(t, err, fmt.Sprintf("Nil for data=%+v", c))
		}
	}
}

func TestGetLogger(t *testing.T) {
	setupLoggerTest(t, configTest)
	r := GetLogger()
	assert.NotNil(t, r)
	logTest(r)
}

func TestNewLogger(t *testing.T) {
	l := newLoggerTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByConfig(t *testing.T) {
	l := newLogrusByConfigTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByInvalidLevelConfig(t *testing.T) {
	l := NewLoggerByConfig(Configuration{Level: Level(99)})
	assert.NotNil(t, l)
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
