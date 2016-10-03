package logger

import (
	"github.com/Sirupsen/logrus"
	"io"
)

func (f Format) toLogrusFormatter() logrus.Formatter {
	switch f {
	case JSON:
		return new(logrus.JSONFormatter)
	default:
		return new(logrus.TextFormatter)
	}
}

func (o Out) toLogrusOut() (io.Writer, error) {
	return getOutput(o)
}

func (n Level) toLogrusLevel() logrus.Level {
	return logrus.Level(n)
}

type logrusLogger struct {
	baseLogger
	logger *logrus.Logger
}

func (l logrusLogger) toLogrusFields(fields ...Field) logrus.Fields {
	logrusFields := make(map[string]interface{})
	for _, v := range fields {
		logrusFields[v.key] = v.val
	}
	return logrusFields
}

func (l logrusLogger) toInterfaceSlice(fields ...Field) []interface{} {
	logrusFields := make([]interface{}, len(fields))
	for i, v := range fields {
		logrusFields[i] = v.val
	}
	return logrusFields
}

func (l logrusLogger) Debug(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Debug(message)
		return
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Debugf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Debug(message)
}

func (l logrusLogger) Info(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Info(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Infof(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Info(message)
}

func (l logrusLogger) Warn(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Warn(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Warnf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Warn(message)
}

func (l logrusLogger) Error(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Error(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Errorf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Error(message)
}

func (l logrusLogger) Panic(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Panic(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Panicf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Panic(message)
}

func (l logrusLogger) Fatal(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Fatal(message)
	}
	if DefaultConfig.Format == LOGRUSFmtfText {
		l.logger.Fatalf(message, l.toInterfaceSlice(fields...)...)
		return
	}
	l.logger.WithFields(l.toLogrusFields(fields...)).Fatal(message)
}

func setupLogrus(loggerConfig *Configuration) error {
	output, err := loggerConfig.Out.toLogrusOut()
	if err != nil {
		return err
	}
	logrus.SetOutput(output)
	logrus.SetFormatter(loggerConfig.Format.toLogrusFormatter())
	logrus.SetLevel(loggerConfig.Level.toLogrusLevel())
	return nil
}

func newLogrus(options ...Option) Logger {
	logger := new(logrusLogger)
	logger.logger = logrus.New()
	return logger
}
