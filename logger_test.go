package l

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	errMock = errors.New("ErrMock")
	objMock = user{
		Name:     "Mock User",
		Email:    "user@mock.com",
		Birthday: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
	}
)

type user struct {
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Birthday time.Time `json:"birthday"`
}

func lFakeFields() []Field {
	return []Field{
		String("string", "string logger field"),
		Bytes("bytes", []byte("[]byte logger field")),
		Int("int", 1),
		Int32("int32", 1),
		Int64("int64", 2),
		Float("float", 3.0),
		Float64("float", 3.0),
		Bool("bool", true),
		Duration("duration", time.Second),
		Time("time", time.Unix(0, 0)),
		Err(errMock),
	}
}

func init() {
	//DefaultConfig
	setupErr := Setup(new(Configuration), DefaultLog, DefaultFieldAdapter())
	if setupErr != nil {
		panic(setupErr)
	}
}

func newLogger(t assert.TestingT) Logger {
	logger, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	return logger
}

func newConfigLogger(t assert.TestingT, c *Configuration) Logger {
	logger, err := NewByConfig(c)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	return logger
}

func logTest(t assert.TestingT, logger Logger) {
	logger.Debug("DebugMessage")
	logger.Debug("DebugFieldsMessage", lFakeFields()...)

	logger.Info("InfoMessage")
	logger.Info("InfoFieldsMessage", lFakeFields()...)

	logger.Warn("WarnMessage")
	logger.Warn("WarnFieldsMessage", lFakeFields()...)

	logger.Error("ErrorMessage")
	logger.Error("ErrorFieldsMessage", lFakeFields()...)

	assert.NotPanics(t, func() {
		logger.Panic("PanicErrorMessage")
	})
	assert.NotPanics(t, func() {
		logger.Panic("PanicFieldsMessage", lFakeFields()...)
	})
}

func TestSetupLogger(t *testing.T) {
	cases := []struct {
		provider     Provider
		fieldAdapter FieldAdapter
		err          error
	}{
		{DefaultLog, DefaultFieldAdapter(), nil},
		{nil, DefaultFieldAdapter(), ErrInvalidProvider},
		{DefaultLog, nil, ErrInvalidFieldAdapter},
		{nil, nil, ErrInvalidProvider},
	}
	for _, c := range cases {
		err := Setup(new(Configuration), c.provider, c.fieldAdapter)
		assert.Equal(t, c.err, err, fmt.Sprintf("Invalid err for data=%+v", c))
	}
}

func TestNewLogger(t *testing.T) {
	l := newLogger(t)
	logTest(t, l)
}

func TestNewConfigLogger(t *testing.T) {
	l := newConfigLogger(t, &Configuration{Format: JSON, Out: STDOUT})
	logTest(t, l)
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
