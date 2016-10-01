package logger

import (
	logrus "github.com/Sirupsen/logrus"
)

//Field is a struct to send paramaters to log messages
type Field struct {
	key string
	val interface{}
}

//Logger is an interface to write log messages
type Logger interface {
	Level() Level
	SetLevel(Level)
	With(...Field) Logger
	IsEnabled(Level, string) bool
	Log(Level, string, ...Field)
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
	Fatal(string, ...Field)
}

type logrusLogger struct {
	logger logrus.Logger
}

func (l *logrusLogger) toLogrusFields(fields ...Field) logrus.Fields {
	var logrusFields logrus.Fields
	for _, v := range fields {
		logrusFields[v.key] = v.val
	}
	return logrusFields
}

func (l *logrusLogger) Log(level Level, message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Info(message)
	} else {

	}

}
