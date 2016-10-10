package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	once          sync.Once
	loggerFactory func(Configuration) Logger
	rootLogger    Logger
)

//Setup initializes the logger system
func Setup(loggerConfig Configuration) error {
	fmt.Printf("logger.Setup config=%+v\n", loggerConfig)
	if loggerConfig.Level == Level(0) {
		loggerConfig.Level = DEBUG
	}
	if loggerConfig.Format == Format("") {
		loggerConfig.Format = TEXT
	}
	if loggerConfig.Out == Out("") {
		loggerConfig.Out = STDOUT
	}
	var setupErr error
	switch loggerConfig.Provider {
	case LOGRUS:
		setupErr = setupLogrus(loggerConfig)
	case ZAP:
		setupErr = setupZap(loggerConfig)
	default:
		setupErr = ErrInvalidProvider
	}
	if setupErr != nil {
		return setupErr
	}
	rootLogger = loggerFactory(loggerConfig)
	DefaultConfig = &loggerConfig
	return nil
}

//Get gets an implemetor of the configured log provider
func Get() Logger {
	// once.Do(func() {
	// if rootLogger == nil {
	// 	rootLogger = NewLoggerByConfig(Configuration{})
	// }
	// })
	if rootLogger == nil {
		return NewLoggerByConfig(Configuration{Provider: LOGRUS})
	}
	return rootLogger
}

//GetLogger gets an implemetor of the configured log provider
func GetLogger() Logger {
	//if rootLogger == nil {
	//	rootLogger = NewLogger()
	//}
	fmt.Printf("logger.GetLogger rootLogger=%+v\n", rootLogger)
	return rootLogger
}

//NewLogger creates a logger implemetor with the default configuration
func NewLogger() Logger {
	fmt.Printf("logger.New config=%+v\n", DefaultConfig)
	return loggerFactory(*DefaultConfig)
}

//NewLoggerByConfig creates a logger implemetor with the provided configuration
func NewLoggerByConfig(config Configuration) Logger {
	switch config.Provider {
	case ZAP:
		return newZap(config)
	default:
		return newLogrus(config)
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

func String(key, val string) Field {
	return Field{key: key, val: val, valType: StringField}
}

func Bytes(key string, val []byte) Field {
	return Field{key: key, val: val, valType: BytesField}
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

func Error(val error) Field {
	return Field{key: "error", val: val, valType: ErrorField}
}
