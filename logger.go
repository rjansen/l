package l

import (
	"errors"
	"fmt"
	"github.com/matryer/resync"
	"time"
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was provided
	ErrInvalidProvider = errors.New("l.ErrInvalidProvider Message='The configured provider is invalid")
	//ErrInvalidSetup is raised when the Setup method does not call or the provider setup is invalid
	ErrInvalidSetup = errors.New("l.ErrInvalidSetup Message='You must call the provider setup function before execute this action'")
	once            resync.Once
	loggerProvider  Provider
	rootLogger      Logger
)

func init() {
	fmt.Printf("l.init\n")
}

//Setup initializes the logger system
func Setup(provider Provider) error {
	if provider == nil {
		return ErrInvalidProvider
	}
	loggerProvider = provider
	rootLogger = loggerProvider()
	return nil
}

func setted() bool {
	return loggerProvider != nil
}

//Get gets an implemetor of the configured log provider
func Get() Logger {
	once.Do(func() {
		if !setted() {
			if err := Setup(newLogger); err != nil {
				panic(err)
			}
		}
	})
	return rootLogger
}

//New creates a logger implemetor with the provided fields
func New(field ...Field) Logger {
	if !setted() {
		panic(ErrInvalidSetup)
	}
	return loggerProvider(field...)
}

func Debug(message string, fields ...Field) {
	Get().Debug(message, fields...)
}

func Info(message string, fields ...Field) {
	Get().Info(message, fields...)

}

func Warn(message string, fields ...Field) {
	Get().Warn(message, fields...)
}

func Error(message string, fields ...Field) {
	Get().Error(message, fields...)
}

func Fatal(message string, fields ...Field) {
	Get().Fatal(message, fields...)
}

func String(key, val string) Field {
	return Field{Key: key, Val: val, ValType: StringField}
}

func Bytes(key string, val []byte) Field {
	return Field{Key: key, Val: string(val), ValType: BytesField}
}

func Int(key string, val int) Field {
	return Field{Key: key, Val: val, ValType: IntField}
}

func Int64(key string, val int64) Field {
	return Field{Key: key, Val: val, ValType: Int64Field}
}

func Float(key string, val float32) Field {
	return Field{Key: key, Val: val, ValType: FloatField}
}

func Float64(key string, val float64) Field {
	return Field{Key: key, Val: val, ValType: Float64Field}
}

func Bool(key string, val bool) Field {
	return Field{Key: key, Val: val, ValType: BoolField}
}

func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Val: val, ValType: DurationField}
}

func Time(key string, val time.Time) Field {
	return Field{Key: key, Val: val, ValType: TimeField}
}

func Struct(key string, val interface{}) Field {
	return Field{Key: key, Val: val, ValType: StructField}
}

func Slice(key string, val interface{}) Field {
	return Field{Key: key, Val: val, ValType: SliceField}
}

func Err(val error) Field {
	return Field{Key: "error", Val: val, ValType: ErrorField}
}
