package l

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger interface {
	Check(zapcore.Level, string) zapWriter
	Sync() error
}

type zapWriter interface {
	Write(...zap.Field)
}

type zapLoggerDelegate struct {
	*zap.Logger
}

func newZapLoggerDelegate(logger *zap.Logger) zapLogger {
	return &zapLoggerDelegate{
		Logger: logger,
	}
}

func (logger *zapLoggerDelegate) Check(level zapcore.Level, msg string) zapWriter {
	return logger.Logger.Check(level, msg)
}

type zapWriterDelegate struct {
	zapWriter
}

func (writer *zapWriterDelegate) Write(values ...Value) {
	fields := make([]zapcore.Field, len(values))
	for index, logValue := range values {
		fields[index] = zap.Any(logValue.name, logValue.value)
	}
	writer.zapWriter.Write(fields...)
}

type zapDriver struct {
	logger zapLogger
}

func (driver zapDriver) Log(level Level, msg string) LogWriter {
	var (
		zapLevel zapcore.Level
		levelErr = zapLevel.Set(level.String())
	)
	if levelErr != nil {
		return nil
	}
	writer := driver.logger.Check(zapLevel, msg)
	if writer == nil {
		return nil
	}
	return &zapWriterDelegate{
		zapWriter: writer,
	}
}

func (driver zapDriver) Close() {
	_ = driver.logger.Sync()
}

func NewDriver(logger zapLogger) zapDriver {
	return zapDriver{
		logger: logger,
	}
}

func NewZapLogger(level Level, output Out) (*zap.Logger, error) {
	var (
		zapLevel  zapcore.Level
		errLevel  = zapLevel.Set(level.String())
		zapOutput = output.String()
	)
	if errLevel != nil {
		return nil, errLevel
	}
	cfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(zapLevel),
		DisableCaller: true,
		Development:   false,
		Encoding:      "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			StacktraceKey:  "stack",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{zapOutput},
		ErrorOutputPaths: []string{zapOutput},
	}
	return cfg.Build()
}

func NewZapLoggerDefault() Logger {
	zapLogger, _ := NewZapLogger(DEBUG, STDOUT)
	return New(
		NewDriver(newZapLoggerDelegate(zapLogger)),
	)
}
