package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

//Setup initializes the logger system
func Setup(loggerConfig *Configuration) error {
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
	case LOGGING:
		setupErr = setupLogging(loggerConfig)
	default:
		setupErr = ErrInvalidProvider
	}
	if setupErr != nil {
		return setupErr
	}
	DefaultConfig = loggerConfig
	return nil
}

//GetLogger gets an implemetor of the configured log provider
func GetLogger(module string) Logger {
	switch DefaultConfig.Provider {
	case LOGRUS:
		return newLogrus()
	case LOGGING:
		return newLogging()
	default:
		return newZap()
	}
}

//NewLogger creates an implemetor of the configured log provider with the provided options + default options
func NewLogger(options ...Option) Logger {
	switch DefaultConfig.Provider {
	case LOGRUS:
		return newLogrus(options...)
	case LOGGING:
		return newLogging(options...)
	default:
		return newZap(options...)
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
	return Field{key: key, val: val}
}

func Int(key string, val int) Field {
	return Field{key: key, val: val}
}

func Float(key string, val float64) Field {
	return Field{key: key, val: val}
}

func Bool(key string, val bool) Field {
	return Field{key: key, val: val}
}

func Struct(key string, val interface{}) Field {
	return Field{key: key, val: val}
}
