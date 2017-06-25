package l

import (
	"errors"
	"time"
)

const (
	NAME_FIELD = "logName"
	LOGDEFAULT = "logDefault"
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was provided
	ErrInvalidProvider = errors.New("l.ErrInvalidProvider Message='The configured provider is invalid")
	//ErrInvalidFieldAdapter is the err raised when an invalid field adapter was provided
	ErrInvalidFieldAdapter = errors.New("l.ErrInvalidFieldAdapter Message='The configured field adapter is invalid")
	//ErrInvalidSetup is raised when the Setup method does not call or the provider setup is invalid
	ErrInvalidSetup = errors.New("l.ErrInvalidSetup Message='You must call the provider setup function before execute this action'")
	defaultConfig   *Configuration
	defaultLog      Logger
	loggerProvider  Provider
	fieldAdapter    FieldAdapter
)

//Setup initializes the logger system
func Setup(c *Configuration, p Provider, a FieldAdapter) error {
	if p == nil {
		return ErrInvalidProvider
	}
	if a == nil {
		return ErrInvalidFieldAdapter
	}
	if c == nil {
		c = new(Configuration)
	}
	defaultConfig = c
	loggerProvider = p
	fieldAdapter = a
	log, err := New(
		String(NAME_FIELD, LOGDEFAULT),
	)
	if err != nil {
		return ErrInvalidProvider
	}
	defaultLog = log
	return nil
}

func Setted() bool {
	return loggerProvider != nil && fieldAdapter != nil
}

//New creates a logger implemetor with the provided fields
func New(fields ...Field) (Logger, error) {
	if !Setted() {
		return nil, ErrInvalidSetup
	}
	return loggerProvider(defaultConfig, fields...)
}

//NewByConfig creates a logger implemetor with the provided fields and configuration
func NewByConfig(cfg *Configuration, fields ...Field) (Logger, error) {
	if !Setted() {
		return nil, ErrInvalidSetup
	}
	return loggerProvider(cfg, fields...)
}

func String(key, val string) Field {
	return fieldAdapter.String(key, val)
}

func Bytes(key string, val []byte) Field {
	return fieldAdapter.Bytes(key, val)
}

func Int(key string, val int) Field {
	return fieldAdapter.Int(key, val)
}

func Int32(key string, val int32) Field {
	return fieldAdapter.Int32(key, val)
}

func Int64(key string, val int64) Field {
	return fieldAdapter.Int64(key, val)
}

func Float(key string, val float32) Field {
	return fieldAdapter.Float(key, val)
}

func Float64(key string, val float64) Field {
	return fieldAdapter.Float64(key, val)
}

func Bool(key string, val bool) Field {
	return fieldAdapter.Bool(key, val)
}

func Duration(key string, val time.Duration) Field {
	return fieldAdapter.Duration(key, val)
}

func Time(key string, val time.Time) Field {
	return fieldAdapter.Time(key, val)
}

func Err(val error) Field {
	return fieldAdapter.Error(val)
}

func Struct(key string, val interface{}) Field {
	return fieldAdapter.Struct(key, val)
}

func Slice(key string, val interface{}) Field {
	return fieldAdapter.Slice(key, val)
}

func DefaultLevel() Level {
	return defaultLog.Level()
}

func Enabled(level Level) bool {
	return defaultLog.Enabled(level)
}

func WithFields(fields ...Field) Logger {
	return defaultLog.WithFields(fields...)
}

func Debug(message string, fields ...Field) {
	defaultLog.Debug(message, fields...)
}

func Info(message string, fields ...Field) {
	defaultLog.Info(message, fields...)
}

func Warn(message string, fields ...Field) {
	defaultLog.Warn(message, fields...)
}

func Error(message string, fields ...Field) {
	defaultLog.Error(message, fields...)
}

func Panic(message string, fields ...Field) {
	defaultLog.Panic(message, fields...)
}

func Fatal(message string, fields ...Field) {
	defaultLog.Fatal(message, fields...)
}

func DefaultString() string {
	return defaultLog.String()
}
