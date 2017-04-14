package logrus

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/rjansen/l"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

const (
	messageMock = "This is a logrus long message mock"
)

var (
	errMock = errors.New("ErrLogrusMock")
	objMock = user{
		Name:     "Logrus Mock User",
		Email:    "logrus.user@mock.com",
		Birthday: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	configTest    = &l.Configuration{Format: l.JSON, Out: l.DISCARD}
	lLogrusConfig = &l.Configuration{Format: l.JSON, Out: l.DISCARD}
)

type user struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

func lFakeFields() []l.Field {
	return []l.Field{
		l.String("string", "string logger field"),
		l.Bytes("bytes", []byte("[]byte logger field")),
		l.Int("int", 1),
		l.Int32("int32", 1),
		l.Int64("int64", 2),
		l.Float("float", 3.0),
		l.Float64("float", 3.0),
		l.Bool("bool", true),
		l.Duration("duration", time.Second),
		l.Time("time", time.Now()),
		l.Err(errMock),
	}
}

func logrusFakeFields() logrus.Fields {
	return logrus.Fields{
		"string":   "string logger field",
		"bytes":    []byte("[]byte logger field"),
		"int":      1,
		"int32":    1,
		"int64":    2,
		"float":    3.0,
		"float64":  float64(3.0),
		"bool":     true,
		"duration": time.Second,
		"time":     time.Now(),
		"error":    errMock,
	}
}

func lSetup(t assert.TestingT, config *l.Configuration) {
	setupErr := Setup(config)
	assert.Nil(t, setupErr)
}

func lNew(t assert.TestingT) l.Logger {
	l, err := l.New()
	assert.Nil(t, err)
	assert.NotNil(t, l)
	return l
}

func lNewByConfig(t assert.TestingT, c *l.Configuration) l.Logger {
	l, err := l.NewByConfig(c)
	assert.Nil(t, err)
	assert.NotNil(t, l)
	return l
}

func lSetupNew(t assert.TestingT, config *l.Configuration) l.Logger {
	lSetup(t, config)
	return lNew(t)
}

func logrusNew(t assert.TestingT) *logrus.Logger {
	l := &logrus.Logger{
		Level:     logrus.DebugLevel,
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Out:       ioutil.Discard,
	}
	assert.NotNil(t, l)
	return l
}

func lLog(t assert.TestingT, logger l.Logger) {
	logger.Debug(messageMock)
	logger.Debug(messageMock, lFakeFields()...)

	logger.Info(messageMock)
	logger.Info(messageMock, lFakeFields()...)

	logger.Warn(messageMock)
	logger.Warn(messageMock, lFakeFields()...)

	logger.Error(messageMock)
	logger.Error(messageMock, lFakeFields()...)

	assert.Panics(t, func() {
		logger.Panic(messageMock)
	})
	assert.Panics(t, func() {
		logger.Panic(messageMock, lFakeFields()...)
	})
}

func logrusLog(t assert.TestingT, logger *logrus.Logger) {
	logger.Debug(messageMock)
	logger.Debug(messageMock, logrusFakeFields())

	logger.Info(messageMock)
	logger.Info(messageMock, logrusFakeFields())

	logger.Warn(messageMock)
	logger.Warn(messageMock, logrusFakeFields())

	logger.Error(messageMock)
	logger.Error(messageMock, logrusFakeFields())

	assert.Panics(t, func() {
		logger.Panic(messageMock)
	})
	assert.Panics(t, func() {
		logger.Panic(messageMock, logrusFakeFields())
	})
}

func TestSetup(t *testing.T) {
	cases := []struct {
		config  l.Configuration
		success bool
	}{
		{*configTest, true},
		{*lLogrusConfig, true},
		{l.Configuration{Debug: true}, true},
		{l.Configuration{Out: l.STDOUT}, true},
		{l.Configuration{Out: l.STDERR}, true},
		{l.Configuration{Out: l.DISCARD}, true},
		{l.Configuration{Out: l.Out("/tmp/glive_test.log")}, true},
		{l.Configuration{Format: l.JSON}, true},
		{l.Configuration{Format: l.TEXT}, true},
		{l.Configuration{Format: l.TEXTColor}, true},
		{l.Configuration{Format: l.JSONColor}, true},
		{l.Configuration{Level: l.Level("invalid")}, false},
	}
	for _, c := range cases {
		err := Setup(&c.config)
		if c.success {
			assert.Nil(t, err, fmt.Sprintf("Not Nil for data=%+v", c))
		} else {
			assert.NotNil(t, err, fmt.Sprintf("Nil for data=%+v", c))
		}
	}
}

func TestNew(t *testing.T) {
	cases := []struct {
		config l.Configuration
		err    error
	}{
		{*configTest, nil},
		{*lLogrusConfig, nil},
		{l.Configuration{Debug: true, Out: l.DISCARD}, nil},
		{l.Configuration{}, nil},
		{l.Configuration{Format: l.TEXT}, nil},
	}
	for _, c := range cases {
		logger := lSetupNew(t, &c.config)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.(*logrusLogger).logger)
		lLog(t, logger)
	}
}

func TestNewByConfig(t *testing.T) {
	cases := []struct {
		config l.Configuration
		err    error
	}{
		{*configTest, nil},
		{*lLogrusConfig, nil},
		{l.Configuration{Out: l.DISCARD}, nil},
		{l.Configuration{}, nil},
	}
	for _, c := range cases {
		logger := lNewByConfig(t, &c.config)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.(*logrusLogger).logger)
		lLog(t, logger)
	}
}

func TestLogrusNew(t *testing.T) {
	logger := logrusNew(t)
	assert.NotNil(t, logger)
	logrusLog(t, logger)
}
