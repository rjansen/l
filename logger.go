package l

import (
	"errors"
	"time"
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was provided
	ErrInvalidProvider = errors.New("l.ErrInvalidProvider Message='The configured provider is invalid")
	//ErrInvalidFieldProvider is the err raised when an invalid field provider was provided
	ErrInvalidFieldProvider = errors.New("l.ErrInvalidFieldProvider Message='The configured field provider is invalid")
	//ErrInvalidSetup is raised when the Setup method does not call or the provider setup is invalid
	ErrInvalidSetup = errors.New("l.ErrInvalidSetup Message='You must call the provider setup function before execute this action'")
	loggerProvider  Provider
	fieldProvider   FieldProvider
)

//Setup initializes the logger system
func Setup(p Provider, fp FieldProvider) error {
	if p == nil {
		return ErrInvalidProvider
	}
	if fp == nil {
		return ErrInvalidFieldProvider
	}
	loggerProvider = p
	fieldProvider = fp
	return nil
}

func setted() bool {
	return loggerProvider != nil
}

//New creates a logger implemetor with the provided fields
func New(field ...Field) (Logger, error) {
	if !setted() {
		return nil, ErrInvalidSetup
	}
	return loggerProvider(field...)
}

func String(key, val string) Field {
	return fieldProvider.String(key, val)
}

func Bytes(key string, val []byte) Field {
	return fieldProvider.Bytes(key, val)
}

func Int(key string, val int) Field {
	return fieldProvider.Int(key, val)
}

func Int32(key string, val int32) Field {
	return fieldProvider.Int32(key, val)
}

func Int64(key string, val int64) Field {
	return fieldProvider.Int64(key, val)
}

func Float(key string, val float32) Field {
	return fieldProvider.Float(key, val)
}

func Float64(key string, val float64) Field {
	return fieldProvider.Float64(key, val)
}

func Bool(key string, val bool) Field {
	return fieldProvider.Bool(key, val)
}

func Duration(key string, val time.Duration) Field {
	return fieldProvider.Duration(key, val)
}

func Time(key string, val time.Time) Field {
	return fieldProvider.Time(key, val)
}

func Struct(key string, val interface{}) Field {
	return fieldProvider.Struct(key, val)
}

func Slice(key string, val interface{}) Field {
	return fieldProvider.Slice(key, val)
}

func Err(val error) Field {
	return fieldProvider.Error(val)
}
