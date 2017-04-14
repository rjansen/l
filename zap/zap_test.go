package zap

import (
	"errors"
	"fmt"
	"github.com/rjansen/l"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// "go.uber.org/zap/zaptest"
	"io/ioutil"
	"testing"
	"time"
)

const (
	messageMock = "This is a long message mock"
)

var (
	errMock = errors.New("ErrMock")
	objMock = user{
		Name:     "Mock User",
		Email:    "user@mock.com",
		Birthday: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
	configTest = &l.Configuration{Format: l.JSON, Out: l.DISCARD}
	lZapConfig = &l.Configuration{Format: l.JSON, Out: l.DISCARD}
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

func zapFakeFields() []zapcore.Field {
	return []zapcore.Field{
		zap.String("string", "string logger field"),
		zap.ByteString("bytes", []byte("[]byte logger field")),
		zap.Int("int", 1),
		zap.Int32("int32", 1),
		zap.Int64("int64", 2),
		zap.Float32("float", 3.0),
		zap.Float64("float", 3.0),
		zap.Bool("bool", true),
		zap.Duration("duration", time.Second),
		zap.Time("time", time.Now()),
		zap.Error(errMock),
	}
}

func sugarFakeFields() []interface{} {
	return []interface{}{
		"string", "string logger field",
		"bytes", "[]byte logger field",
		"int", 1,
		"int32", 1,
		"int64", 2,
		"float", 3.0,
		"float", 3.0,
		"bool", true,
		"duration", time.Second,
		"time", time.Now(),
		errMock,
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

func zapNew(t assert.TestingT) *zap.Logger {
	l := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
			zapcore.AddSync(ioutil.Discard),
			zap.DebugLevel,
		),
	)
	assert.NotNil(t, l)
	return l
}

func sugarNew(t assert.TestingT) *zap.SugaredLogger {
	l := zapNew(t)
	return l.Sugar()
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

func zapLog(t assert.TestingT, logger *zap.Logger) {
	logger.Debug(messageMock)
	logger.Debug(messageMock, zapFakeFields()...)

	logger.Info(messageMock)
	logger.Info(messageMock, zapFakeFields()...)

	logger.Warn(messageMock)
	logger.Warn(messageMock, zapFakeFields()...)

	logger.Error(messageMock)
	logger.Error(messageMock, zapFakeFields()...)

	assert.Panics(t, func() {
		logger.Panic(messageMock)
	})
	assert.Panics(t, func() {
		logger.Panic(messageMock, zapFakeFields()...)
	})
}

func sugarLog(t assert.TestingT, logger *zap.SugaredLogger) {
	logger.Debugw(messageMock)
	logger.Debugw(messageMock, sugarFakeFields()...)

	logger.Infow(messageMock)
	logger.Infow(messageMock, sugarFakeFields()...)

	logger.Warnw(messageMock)
	logger.Warnw(messageMock, sugarFakeFields()...)

	logger.Errorw(messageMock)
	logger.Errorw(messageMock, sugarFakeFields()...)

	assert.Panics(t, func() {
		logger.Panicw(messageMock)
	})
	assert.Panics(t, func() {
		logger.Panicw(messageMock, sugarFakeFields()...)
	})
}

func TestSetup(t *testing.T) {
	cases := []struct {
		config  l.Configuration
		success bool
	}{
		{*configTest, true},
		{*lZapConfig, true},
		{l.Configuration{Debug: true}, true},
		{l.Configuration{Out: l.STDOUT}, true},
		{l.Configuration{Out: l.STDERR}, true},
		{l.Configuration{Out: l.DISCARD}, true},
		{l.Configuration{Out: l.Out("/tmp/glive_test.log")}, true},
		{l.Configuration{Format: l.JSON}, true},
		{l.Configuration{Format: l.TEXT}, true},
		{l.Configuration{Format: l.TEXTColor}, true},
		{l.Configuration{Format: l.JSONColor}, true},
		{l.Configuration{Level: l.Level("invalid")}, true},
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
		{*lZapConfig, nil},
		{l.Configuration{Debug: true, Out: l.DISCARD}, nil},
		{l.Configuration{}, nil},
		{l.Configuration{Format: l.TEXT}, nil},
	}
	for _, c := range cases {
		logger := lSetupNew(t, &c.config)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.(*zapLogger).logger)
		lLog(t, logger)
	}
}

func TestNewByConfig(t *testing.T) {
	cases := []struct {
		config l.Configuration
		err    error
	}{
		{*configTest, nil},
		{*lZapConfig, nil},
		{l.Configuration{Out: l.DISCARD}, nil},
		{l.Configuration{}, nil},
	}
	for _, c := range cases {
		logger := lNewByConfig(t, &c.config)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.(*zapLogger).logger)
		lLog(t, logger)
	}
}

func TestZapNew(t *testing.T) {
	logger := zapNew(t)
	assert.NotNil(t, logger)
	zapLog(t, logger)
}

func TestSugarNew(t *testing.T) {
	logger := sugarNew(t)
	assert.NotNil(t, logger)
	sugarLog(t, logger)
}
