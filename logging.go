package logger

import (
	"github.com/op/go-logging"
)

var (
	defaultLoggingFormat = "%{time:2006-01-02T15:04:05.999Z-07:00} %{id:03x} [%{level:.5s}] %{shortpkg}.%{longfunc} %{message}"
	loggingFormatter     = logging.MustStringFormatter(defaultLoggingFormat)
)

type loggingLogger struct {
	baseLogger
	logger *logging.Logger
}

func (l loggingLogger) toLoggingFields(fields ...Field) []interface{} {
	var loggingFields []interface{}
	for _, v := range fields {
		loggingFields = append(loggingFields, v.val)
	}
	return loggingFields
}

func (l loggingLogger) Debug(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Debugf(message, l.toLoggingFields(fields...)...)
	} else {
		l.logger.Debug(message)
	}
}

func (l loggingLogger) Info(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Infof(message, l.toLoggingFields(fields...)...)
	} else {
		l.logger.Info(message)
	}
}

func (l loggingLogger) Warn(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Warningf(message, l.toLoggingFields(fields...)...)
	} else {
		l.logger.Warning(message)
	}
}

func (l loggingLogger) Error(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Errorf(message, l.toLoggingFields(fields...)...)
	} else {
		l.logger.Error(message)
	}
}

func (l loggingLogger) Panic(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Panicf(message, l.toLoggingFields(fields...)...)
	} else {
		l.logger.Panic(message)
	}
}

func (l loggingLogger) Fatal(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Fatalf(message, l.toLoggingFields(fields...)...)
	} else {
		l.logger.Fatal(message)
	}
}

func setupLogging(loggerConfig *Configuration) error {
	output, err := getOutput(loggerConfig.Out)
	if err != nil {
		return err
	}
	backEndMessages := logging.NewBackendFormatter(logging.NewLogBackend(output, "", 0), loggingFormatter)
	levelMessages := logging.AddModuleLevel(backEndMessages)
	levelMessages.SetLevel(logging.Level(loggerConfig.Level), "")
	logging.SetBackend(levelMessages)
	return nil
}

func newLogging(options ...Option) Logger {
	logger := new(loggingLogger)
	logger.logger = logging.MustGetLogger("rootLogger")
	return logger
}
