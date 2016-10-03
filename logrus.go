package logger

import (
	"github.com/Sirupsen/logrus"
)

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

func (l *logrusLogger) applyOptions(options ...Option) {
	l.logger = logrus.New()
	for _, o := range options {
		o.apply(l)
	}
}

func (l logrusLogger) Debug(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Debug(message)
	} else {
		l.logger.Debug(message)
	}
}

func (l logrusLogger) Info(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Info(message)
	} else {
		l.logger.Info(message)
	}
}

func (l logrusLogger) Warn(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Warn(message)
	} else {
		l.logger.Warn(message)
	}
}

func (l logrusLogger) Error(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Error(message)
	} else {
		l.logger.Error(message)
	}
}

func (l logrusLogger) Panic(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Panic(message)
	} else {
		l.logger.Panic(message)
	}
}

func (l logrusLogger) Fatal(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.WithFields(l.toLogrusFields(fields...)).Fatal(message)
	} else {
		l.logger.Fatal(message)
	}
}

func setupLogrus(loggerConfig *Configuration) error {
	output, err := getOutput(loggerConfig.Out)
	if err != nil {
		return err
	}
	logrus.SetOutput(output)
	defaultOptions = append(defaultOptions, loggerConfig.Out)
	switch loggerConfig.Format {
	case JSON:
		logrus.SetFormatter(new(logrus.TextFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
	defaultOptions = append(defaultOptions, loggerConfig.Format)
	logrus.SetLevel(logrus.Level(loggerConfig.Level))
	defaultOptions = append(defaultOptions, loggerConfig.Level)
	return nil
}

func newLogrus(options ...Option) Logger {
	logger := new(logrusLogger)
	logger.applyOptions(append(defaultOptions, options...)...)
	return logger
}
