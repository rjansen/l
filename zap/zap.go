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
)

func newFieldProvider() *zapFieldProvider {
	return new(zapFieldProvider)
}

type zapField struct {
	zapcore.Field
}

func (f zapField) Key() string {
	return f.Field.Key
}

func (f zapField) Val() interface{} {
	//TODO: Vaidate this approach
	if f.Interface != nil {
		return f.Interface
	}
	if f.Integer > 0 {
		return f.Integer
	}
	return f.String
}

func (f zapField) Type() l.FieldType {
	return l.FieldType(f.Field.Type)
}

func newZapField(f zapcore.Field) *zapField {
	return &zapField{Field: f}
}

type zapFieldProvider struct {
}

func (zapFieldProvider) String(key string, val string) l.Field {
	return newZapField(zap.String(key, val))
}

func (zapFieldProvider) Bytes(key string, val []byte) l.Field {
	return newZapField(zap.Binary(key, val))
}

func (zapFieldProvider) Int(key string, val int) l.Field {
	return newZapField(zap.Int(key, val))
}

func (zapFieldProvider) Int32(key string, val int32) l.Field {
	return newZapField(zap.Int32(key, val))
}

func (zapFieldProvider) Int64(key string, val int64) l.Field {
	return newZapField(zap.Int64(key, val))
}

func (zapFieldProvider) Float(key string, val float32) l.Field {
	return newZapField(zap.Float32(key, val))
}

func (zapFieldProvider) Float64(key string, val float64) l.Field {
	return newZapField(zap.Float64(key, val))
}

func (zapFieldProvider) Duration(key string, val time.Duration) l.Field {
	return newZapField(zap.Duration(key, val))
}

func (zapFieldProvider) Time(key string, val time.Time) l.Field {
	return newZapField(zap.Time(key, val))
}

func (zapFieldProvider) Bool(key string, val bool) l.Field {
	return newZapField(zap.Bool(key, val))
}

func (zapFieldProvider) Struct(key string, val interface{}) l.Field {
	return newZapField(zap.Any(key, val))
}

func (zapFieldProvider) Slice(key string, val interface{}) l.Field {
	return newZapField(zap.Any(key, val))
}

func (zapFieldProvider) Error(val error) l.Field {
	return newZapField(zap.Error(val))
}

func toZapFields(fields ...l.Field) []zapcore.Field {
	var zapFields []zapcore.Field
	for _, v := range fields {
		zapFields = append(zapFields, v.(*zapField).Field)
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

func (l zapLogger) Debug(message string, fields ...l.Field) {
	if len(fields) > 0 {
		l.logger.Debug(message, toZapFields(fields...)...)
	} else {
		l.logger.Debug(message)
	}
}

func (l zapLogger) Info(message string, fields ...l.Field) {
	if len(fields) > 0 {
		l.logger.Info(message, toZapFields(fields...)...)
	} else {
		l.logger.Info(message)
	}
}

func (l zapLogger) Warn(message string, fields ...l.Field) {
	if len(fields) > 0 {
		l.logger.Info(message, toZapFields(fields...)...)
	} else {
		l.logger.Warn(message)
	}
}

func (l zapLogger) Error(message string, fields ...l.Field) {
	if len(fields) > 0 {
		l.logger.Error(message, toZapFields(fields...)...)
	} else {
		l.logger.Error(message)
	}
}

func (l zapLogger) Panic(message string, fields ...l.Field) {
	if len(fields) > 0 {
		l.logger.Panic(message, toZapFields(fields...)...)
	} else {
		l.logger.Panic(message)
	}
}

func (l zapLogger) Fatal(message string, fields ...l.Field) {
	if len(fields) > 0 {
		l.logger.Fatal(message, toZapFields(fields...)...)
	} else {
		l.logger.Fatal(message)
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
	defaultConfig = zapConfig
	return l.Setup(New, newFieldProvider())
}

func New(field ...l.Field) (l.Logger, error) {
	return create(defaultConfig, field...)
}

func NewByConfig(loggerConfig *l.Configuration, field ...l.Field) (l.Logger, error) {
	zapConfig, errs := toZapConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		if loggerConfig.Debug {
			fmt.Printf("l.zap.NewByConfigErr Config=%s errs=%s\n", zapConfig.String(), errs)
		}
		return nil, fmt.Errorf("l.zap.NewByConfigErr Config=%s Errs=%v", zapConfig.String(), errs)
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
	// logger.logger.Data = toLogrusFields(field...)
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
		level:     level,
		formatter: encoder,
		output:    output,
	}, errs
}
