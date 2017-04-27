package zap

import (
	"fmt"
	"github.com/rjansen/l"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"time"
)

var (
	defaultConfig *zapConfig
	lConfig       *l.Configuration
)

func newFieldAdapter() *ZapFieldAdapter {
	return new(ZapFieldAdapter)
}

// type LazyZapField struct {
// 	zapcore.Field
// 	Producer func() zapcore.Field
// }

// func (LazyZapField) Key() {

// }

type ZapField zapcore.Field

func (f ZapField) Name() string {
	return f.Key
}

func (f ZapField) Value() interface{} {
	//TODO: Vaidate this approach
	if f.Interface != nil {
		return f.Interface
	}
	if f.Integer > 0 {
		return f.Integer
	}
	return f.String
}

func (f ZapField) ZapField() zapcore.Field {
	return zapcore.Field(f)
}

type ZapFieldAdapter struct {
}

func (ZapFieldAdapter) String(key string, val string) l.Field {
	return ZapField(zap.String(key, val))
}

func (ZapFieldAdapter) Bytes(key string, val []byte) l.Field {
	return ZapField(zap.Binary(key, val))
}

func (ZapFieldAdapter) Int(key string, val int) l.Field {
	return ZapField(zap.Int(key, val))
}

func (ZapFieldAdapter) Int32(key string, val int32) l.Field {
	return ZapField(zap.Int32(key, val))
}

func (ZapFieldAdapter) Int64(key string, val int64) l.Field {
	return ZapField(zap.Int64(key, val))
}

func (ZapFieldAdapter) Float(key string, val float32) l.Field {
	return ZapField(zap.Float32(key, val))
}

func (ZapFieldAdapter) Float64(key string, val float64) l.Field {
	return ZapField(zap.Float64(key, val))
}

func (ZapFieldAdapter) Duration(key string, val time.Duration) l.Field {
	return ZapField(zap.Duration(key, val))
}

func (ZapFieldAdapter) Time(key string, val time.Time) l.Field {
	return ZapField(zap.Time(key, val))
}

func (ZapFieldAdapter) Bool(key string, val bool) l.Field {
	return ZapField(zap.Bool(key, val))
}

func (ZapFieldAdapter) Error(val error) l.Field {
	return ZapField(zap.Error(val))
}

func (ZapFieldAdapter) Struct(key string, val interface{}) l.Field {
	return ZapField(zap.Any(key, val))
}

func (ZapFieldAdapter) Slice(key string, val interface{}) l.Field {
	return ZapField(zap.Any(key, val))
}

func toZapFields(fields ...l.Field) []zapcore.Field {
	fieldsLen := len(fields)
	if fieldsLen <= 0 {
		return nil
	}
	zapFields := make([]zapcore.Field, fieldsLen)
	// var zapFields []zapcore.Field
	for i, v := range fields {
		zapFields[i] = v.(ZapField).ZapField()
		// zapFields = append(zapFields, field)
	}
	return zapFields
}

type zapLogger struct {
	l.BaseLogger
	logger *zap.Logger
}

func (l *zapLogger) WithFields(fields ...l.Field) l.Logger {
	return &zapLogger{
		logger: l.logger.With(toZapFields(fields...)...),
	}
}

func (l *zapLogger) Debug(message string, fields ...l.Field) {
	if ce := l.logger.Check(zap.DebugLevel, message); ce != nil {
		ce.Write(toZapFields(fields...)...)
	}
}

func (l *zapLogger) Info(message string, fields ...l.Field) {
	if ce := l.logger.Check(zap.InfoLevel, message); ce != nil {
		ce.Write(toZapFields(fields...)...)
	}
}

func (l *zapLogger) Warn(message string, fields ...l.Field) {
	if ce := l.logger.Check(zap.WarnLevel, message); ce != nil {
		ce.Write(toZapFields(fields...)...)
	}
}

func (l *zapLogger) Error(message string, fields ...l.Field) {
	if ce := l.logger.Check(zap.ErrorLevel, message); ce != nil {
		ce.Write(toZapFields(fields...)...)
	}
}

func (l *zapLogger) Panic(message string, fields ...l.Field) {
	if ce := l.logger.Check(zap.PanicLevel, message); ce != nil {
		ce.Write(toZapFields(fields...)...)
	}
}

func (l *zapLogger) Fatal(message string, fields ...l.Field) {
	if ce := l.logger.Check(zap.FatalLevel, message); ce != nil {
		ce.Write(toZapFields(fields...)...)
	}
}

func (zapLogger) String() string {
	return "provider=zap"
}

func Setup(loggerConfig *l.Configuration) error {
	zapConfig, errs := toZapConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		if loggerConfig.Debug {
			fmt.Printf("l.zap.SetupErr Config=%s Errs=%v\n", zapConfig.String(), errs)
		}
		return fmt.Errorf("l.zap.SetupErr Config=%s Errs=%v", zapConfig.String(), errs)
	}
	lConfig = loggerConfig
	defaultConfig = zapConfig
	return l.Setup(lConfig, New, newFieldAdapter())
}

func New(loggerConfig *l.Configuration, field ...l.Field) (l.Logger, error) {
	if lConfig == loggerConfig {
		return create(defaultConfig, field...)
	}
	zapConfig, errs := toZapConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		if loggerConfig.Debug {
			fmt.Printf("l.zap.NewErr Config=%s errs=%s\n", zapConfig.String(), errs)
		}
		return nil, fmt.Errorf("l.zap.NewErr Config=%s Errs=%v", zapConfig.String(), errs)
	}
	return create(zapConfig, field...)
}

func create(cfg *zapConfig, field ...l.Field) (l.Logger, error) {
	if cfg.debug {
		fmt.Printf("l.zap.Creating Config=%s\n", cfg.String())
	}
	logBackend := zap.New(
		zapcore.NewCore(
			cfg.formatter,
			cfg.output,
			cfg.level,
		),
	)
	// for _, hook := range cfg.hooks {
	// 	logBackend.Hooks.Add(hook)
	// }
	logger := new(zapLogger)
	logger.logger = logBackend
	if cfg.debug {
		fmt.Printf("l.zap.Created Config=%s Logger=%s\n", cfg.String(), logger.String())
	}
	return logger, nil
}

type zapConfig struct {
	debug     bool
	output    zapcore.WriteSyncer
	formatter zapcore.Encoder
	level     zapcore.Level
}

func (l zapConfig) String() string {
	return fmt.Sprintf("debug=%t level=%s hasFormatter=%t hasOutput=%t", l.debug, l.level.String(), l.formatter != nil, l.output != nil)
}

func toZapConfig(cfg *l.Configuration) (*zapConfig, []error) {
	var errs []error
	var output zapcore.WriteSyncer
	switch cfg.Out {
	case l.STDOUT, l.Out(""):
		output = zapcore.AddSync(os.Stdout)
	case l.STDERR:
		output = zapcore.AddSync(os.Stderr)
	case l.DISCARD:
		output = zapcore.AddSync(ioutil.Discard)
		// output =  &zaptest.Discarder{}
	default:
		fileOutput, err := cfg.Out.GetOutput()
		//TODO: Think better
		if err != nil {
			panic(err)
		}
		output = zapcore.AddSync(fileOutput)
	}

	var level zapcore.Level
	switch cfg.Level {
	case l.DEBUG:
		level = zap.DebugLevel
	case l.INFO:
		level = zap.InfoLevel
	case l.WARN:
		level = zap.WarnLevel
	case l.ERROR:
		level = zap.ErrorLevel
	case l.PANIC:
		level = zap.PanicLevel
	case l.FATAL:
		level = zap.FatalLevel
	default:
		level = zap.DebugLevel
	}

	var encoder zapcore.Encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	switch cfg.Format {
	case l.TEXT:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	return &zapConfig{
		debug:     cfg.Debug,
		level:     level,
		formatter: encoder,
		output:    output,
	}, errs
}
