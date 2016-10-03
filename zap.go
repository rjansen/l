package logger

import (
	"github.com/uber-go/zap"
	"os"
	"time"
)

var (
	zapFactory zap.Logger
)

func (f Format) toZapEncoder() zap.Encoder {
	switch f {
	case JSON:
		return zap.NewJSONEncoder()
	default:
		return zap.NewTextEncoder()
	}
}

func (o Out) toZapOut() zap.Option {
	switch o {
	case STDOUT:
		return zap.Output(os.Stdout)
	case STDERR:
		return zap.Output(os.Stderr)
	case DISCARD:
		return zap.DiscardOutput
	default:
		fileOutput, _ := getOutput(o)
		zapOutput := zap.AddSync(fileOutput)
		return zap.Output(zapOutput)
	}
}

func (n Level) toZapLevel() zap.Level {
	switch n {
	case DEBUG:
		return zap.DebugLevel
	case INFO:
		return zap.InfoLevel
	case WARN:
		return zap.WarnLevel
	case ERROR:
		return zap.ErrorLevel
	case PANIC:
		return zap.PanicLevel
	case FATAL:
		return zap.FatalLevel
	default:
		return zap.Level(n)
	}
}

type zapLogger struct {
	baseLogger
	logger zap.Logger
}

func (l zapLogger) toZapFields(fields ...Field) []zap.Field {
	var zapFields []zap.Field
	for _, v := range fields {
		switch v.valType {
		case IntField:
			zapFields = append(zapFields, zap.Int64(v.key, v.val.(int64)))
		case StringField:
			zapFields = append(zapFields, zap.String(v.key, v.val.(string)))
		case BoolField:
			zapFields = append(zapFields, zap.Bool(v.key, v.val.(bool)))
		case FloatField:
			zapFields = append(zapFields, zap.Float64(v.key, v.val.(float64)))
		case DurationField:
			zapFields = append(zapFields, zap.Duration(v.key, v.val.(time.Duration)))
		case TimeField:
			zapFields = append(zapFields, zap.Time(v.key, v.val.(time.Time)))
		default:
			zapFields = append(zapFields, zap.Object(v.key, v.val))
		}
	}
	return zapFields
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
	zapFactory = zap.New(loggerConfig.Format.toZapEncoder(), loggerConfig.Level.toZapLevel(), loggerConfig.Out.toZapOut())
	return nil
}

func newZap(options ...Option) Logger {
	logger := new(zapLogger)
	logger.logger = zapFactory.With()
	return logger
}
