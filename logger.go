package logger

import (
	"fmt"
	"github.com/matryer/resync"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	once          resync.Once
	loggerFactory func(Configuration) Logger
)

//Setup initializes the logger system
func Setup(loggerConfig Configuration) error {
	if loggerConfig.Debug {
		fmt.Printf("logger.Setup config=%+v\n", loggerConfig)
	}
	var setupErr error
	switch loggerConfig.Provider {
	case ZAP:
		setupErr = setupZap(loggerConfig)
	default:
		setupErr = setupLogrus(loggerConfig)
	}
	if setupErr != nil {
		return setupErr
	}
	return nil
}

//Get gets an implemetor of the configured log provider
func Get() Logger {
	once.Do(func() {
		if rootLogger == nil {
			setted := isSetted()
			fmt.Println("logger.Get.isSetted =", setted, " defaultConfig =", defaultConfig, " loggerFactory =", loggerFactory)
			if !isSetted() {
				cfg := Configuration{
					Provider: LOGRUS,
					Format:   TEXTColor,
					Out:      STDOUT,
					Hooks:    Hooks("syslog"),
				}
				fmt.Println("logger.Get.Setup =", cfg.String())
				err := Setup(cfg)
				if err != nil {
					panic(err)
				}
			}
			rootLogger = create()
		}
	})
	return rootLogger
}

func isSetted() bool {
	return loggerFactory != nil && defaultConfig != nil
}

func create() Logger {
	setted := isSetted()
	fmt.Println("logger.create.isSetted =", setted, " defaultConfig =", defaultConfig, " loggerFactory =", loggerFactory)
	if !setted {
		panic(ErrSetupNeverCalled)
	}
	return loggerFactory(*defaultConfig)
}

//New creates a logger implemetor with the provided configuration
func New(config Configuration) Logger {
	switch config.Provider {
	case ZAP:
		return newZap(config)
	case LOGRUS:
		return newLogrus(config)
	default:
		panic(ErrInvalidProvider)
	}
}

func getOutput(out Out) (io.Writer, error) {
	switch out {
	case STDOUT:
		return os.Stdout, nil
	case STDERR:
		return os.Stderr, nil
	case DISCARD:
		return ioutil.Discard, nil
	default:
		file, err := os.OpenFile(out.String(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("CreateFileOutputErr[Out=%v Message='%v']", out, err)
		}
		return file, nil
	}
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

func Panic(message string, fields ...Field) {
	Get().Panic(message, fields...)
}

func Fatal(message string, fields ...Field) {
	Get().Fatal(message, fields...)
}

func Debugf(message string, fields ...interface{}) {
	Get().Debugf(message, fields...)
}

func Infof(message string, fields ...interface{}) {
	Get().Infof(message, fields...)
}

func Warnf(message string, fields ...interface{}) {
	Get().Warnf(message, fields...)
}

func Errorf(message string, fields ...interface{}) {
	Get().Errorf(message, fields...)
}

func Panicf(message string, fields ...interface{}) {
	Get().Panicf(message, fields...)
}

func Fatalf(message string, fields ...interface{}) {
	Get().Fatalf(message, fields...)
}

func String(key, val string) Field {
	return Field{key: key, val: val, valType: StringField}
}

func Bytes(key string, val []byte) Field {
	return Field{key: key, val: string(val), valType: BytesField}
}

func Int(key string, val int) Field {
	return Field{key: key, val: val, valType: IntField}
}

func Int64(key string, val int64) Field {
	return Field{key: key, val: val, valType: Int64Field}
}

func Float(key string, val float32) Field {
	return Field{key: key, val: val, valType: FloatField}
}

func Float64(key string, val float64) Field {
	return Field{key: key, val: val, valType: Float64Field}
}

func Bool(key string, val bool) Field {
	return Field{key: key, val: val, valType: BoolField}
}

func Duration(key string, val time.Duration) Field {
	return Field{key: key, val: val, valType: DurationField}
}

func Time(key string, val time.Time) Field {
	return Field{key: key, val: val, valType: TimeField}
}

func Struct(key string, val interface{}) Field {
	return Field{key: key, val: val, valType: StructField}
}

func Err(val error) Field {
	return Field{key: "error", val: val, valType: ErrorField}
}
