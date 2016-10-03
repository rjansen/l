package logger

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/op/go-logging"
	"github.com/uber-go/zap"
	"os"
)

const (
	//STDOUT any message to stdout
	STDOUT Out = "stdout"
	//STDERR redirects any message to stderr
	STDERR = "stderr"
	//DISCARD set logger to ignore all message
	DISCARD = "discard"

	//LOGRUS is the github.com/Sirupsen/logrus id
	LOGRUS Provider = "logrus"
	//ZAP is the github.com/uber-go/zap id
	ZAP = "zap"
	//LOGGING is the github.com/op/go-logging id
	LOGGING = "logging"

	//TEXT is the text log format
	TEXT Format = "text"
	//JSON is the json log format
	JSON = "json"

	//PANIC is the panic level logger
	PANIC Level = iota
	//FATAL is the fatal level logger
	FATAL
	//ERROR is the error level logger
	ERROR
	//WARN is the warn level logger
	WARN
	//INFO is the info level logger
	INFO
	//DEBUG is the debug level logger
	DEBUG
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was select
	ErrInvalidProvider = errors.New("Logger.InvalidProvider[Message='Avaible providers are: LOGRUS, ZAP and LOGGING']")
	//DefaultConfig holds the instance of the behavior parameters
	DefaultConfig  *Configuration
	defaultOptions []Option
)

//Provider is the back end implementor id of the logging feature
type Provider string

//Out is the type for logger writer config
type Out string

func (o Out) String() string {
	return string(o)
}

//Output creates a option for log output
func (o Out) apply(l Logger) error {
	output, err := getOutput(o)
	if err != nil {
		return err
	}
	//TODO: Refactor
	switch DefaultConfig.Provider {
	case LOGRUS:
		logrusLogger := l.(*logrusLogger)
		logrusLogger.logger.Out = output
	case ZAP:
		zapLogger := l.(*zapLogger)
		zapLogger.output = zap.AddSync(output)
	case LOGGING:
		loggingLogger := l.(*loggingLogger)
		loggingLogger.output = output
		backEnd := logging.NewBackendFormatter(logging.NewLogBackend(loggingLogger.output, "", 0), loggingFormatter)
		backEndLeveled := logging.AddModuleLevel(backEnd)
		backEndLeveled.SetLevel(logging.Level(loggingLogger.level), "")
		loggingLogger.logger.SetBackend(backEndLeveled)
	}
	return nil
}

// Option is used to set options for the logger.
type Option interface {
	apply(Logger) error
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(Logger) error

func (f optionFunc) apply(l Logger) error {
	return f(l)
}

//Level is the threshold of the logger
type Level int

func (n Level) apply(l Logger) error {
	//TODO: Refactor
	switch DefaultConfig.Provider {
	case LOGRUS:
		logrusLogger := l.(*logrusLogger)
		logrusLogger.logger.Level = logrus.Level(n)
	case ZAP:
		zapLogger := l.(*zapLogger)
		zapLogger.level = n
		switch n {
		case DEBUG:
			zapLogger.zapLevel = zap.DebugLevel
		case INFO:
			zapLogger.zapLevel = zap.InfoLevel
		case WARN:
			zapLogger.zapLevel = zap.WarnLevel
		case ERROR:
			zapLogger.zapLevel = zap.ErrorLevel
		case PANIC:
			zapLogger.zapLevel = zap.PanicLevel
		case FATAL:
			zapLogger.zapLevel = zap.FatalLevel
		}
	case LOGGING:
		loggingLogger := l.(*loggingLogger)
		loggingLogger.level = n
		backEnd := logging.NewBackendFormatter(logging.NewLogBackend(loggingLogger.output, "", 0), loggingFormatter)
		backEndLeveled := logging.AddModuleLevel(backEnd)
		backEndLeveled.SetLevel(logging.Level(n), "")
		loggingLogger.logger.SetBackend(backEndLeveled)
	}
	return nil
}

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case WARN:
		return "warn"
	case ERROR:
		return "error"
	case PANIC:
		return "panic"
	case FATAL:
		return "fatal"
	default:
		return fmt.Sprintf("Level(%d)", l)
	}
}

//Format is a parameter to controle the logger style
type Format string

func (f Format) apply(l Logger) error {
	//TODO: Refactor
	switch DefaultConfig.Provider {
	case LOGRUS:
		logrusLogger := l.(*logrusLogger)
		switch f {
		case JSON:
			logrusLogger.logger.Formatter = new(logrus.JSONFormatter)
		default:
			logrusLogger.logger.Formatter = new(logrus.TextFormatter)
		}
	case ZAP:
		zapLogger := l.(*zapLogger)
		switch f {
		case JSON:
			zapLogger.encoder = zap.NewJSONEncoder()
		default:
			zapLogger.encoder = zap.NewTextEncoder()
		}
	}
	return nil
}

//Configuration holds the log beahvior parameters
type Configuration struct {
	Provider Provider
	Level    Level
	Format   Format
	Out      Out
}

func (l *Configuration) String() string {
	return fmt.Sprintf("Configuration[Provider=%v Level=%v Format=%v Out=%v]", l.Provider, l.Level, l.Format, l.Out)
}

//FileOutput creates a option for file output
func FileOutput(output *os.File) Option {
	//TODO: Refactor
	return optionFunc(func(l Logger) error {
		switch DefaultConfig.Provider {
		case LOGRUS:
			logrusLogger := l.(*logrusLogger)
			logrusLogger.logger.Out = output
		case ZAP:
			zapLogger := l.(*zapLogger)
			zapLogger.output = output
		case LOGGING:
			loggingLogger := l.(*loggingLogger)
			loggingLogger.output = output
			backEnd := logging.NewBackendFormatter(logging.NewLogBackend(loggingLogger.output, "", 0), loggingFormatter)
			backEndLeveled := logging.AddModuleLevel(backEnd)
			backEndLeveled.SetLevel(logging.Level(loggingLogger.level), "")
			loggingLogger.logger.SetBackend(backEndLeveled)
		}
		return nil
	})
}

//Field is a struct to send paramaters to log messages
type Field struct {
	key string
	val interface{}
}

//Logger is an interface to write log messages
type Logger interface {
	Level() Level
	IsEnabled(Level) bool
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
	Fatal(string, ...Field)
}

type baseLogger struct {
	level Level
}

func (b baseLogger) Level() Level {
	return b.level
}

func (b baseLogger) IsEnabled(level Level) bool {
	return b.level >= level
}
