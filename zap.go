package logger

import (
	"github.com/uber-go/zap"
)

type zapLogger struct {
	baseLogger
	logger   zap.Logger
	zapLevel zap.Level
	//TODO: Refactor
	encoder zap.Encoder
	//TODO: Refactor
	output zap.WriteSyncer
}

func (l zapLogger) toZapFields(fields ...Field) []zap.Field {
	var zapFields []zap.Field
	for _, v := range fields {
		switch v.val.(type) {
		case int:
			zapFields = append(zapFields, zap.Int(v.key, v.val.(int)))
		case int64:
			zapFields = append(zapFields, zap.Int64(v.key, v.val.(int64)))
		case string:
			zapFields = append(zapFields, zap.String(v.key, v.val.(string)))
		case bool:
			zapFields = append(zapFields, zap.Bool(v.key, v.val.(bool)))
		case float64:
			zapFields = append(zapFields, zap.Float64(v.key, v.val.(float64)))
		default:
			zapFields = append(zapFields, zap.Object(v.key, v.val))
		}
	}
	return zapFields
}

func (l *zapLogger) applyOptions(options ...Option) {
	for _, o := range options {
		o.apply(l)
	}
	//TODO: Refactor
	l.logger = zap.New(l.encoder, l.zapLevel, zap.Output(l.output))
}

func (l zapLogger) Debug(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Debug(message, l.toZapFields(fields...)...)
	} else {
		l.logger.Debug(message)
	}
}

func (l zapLogger) Info(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Info(message, l.toZapFields(fields...)...)
	} else {
		l.logger.Info(message)
	}
}

func (l zapLogger) Warn(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Info(message, l.toZapFields(fields...)...)
	} else {
		l.logger.Warn(message)
	}
}

func (l zapLogger) Error(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Error(message, l.toZapFields(fields...)...)
	} else {
		l.logger.Error(message)
	}
}

func (l zapLogger) Panic(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Panic(message, l.toZapFields(fields...)...)
	} else {
		l.logger.Panic(message)
	}
}

func (l zapLogger) Fatal(message string, fields ...Field) {
	if len(fields) > 0 {
		l.logger.Fatal(message, l.toZapFields(fields...)...)
	} else {
		l.logger.Fatal(message)
	}
}

func setupZap(loggerConfig *Configuration) error {
	defaultOptions = append(defaultOptions, loggerConfig.Out)
	defaultOptions = append(defaultOptions, loggerConfig.Format)
	defaultOptions = append(defaultOptions, loggerConfig.Level)
	return nil
}

func newZap(options ...Option) Logger {
	logger := new(zapLogger)
	logger.applyOptions(append(defaultOptions, options...)...)
	return logger
}
