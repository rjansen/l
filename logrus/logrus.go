package logrus

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"time"
	// logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/rjansen/l"
	"io"
	// "log/syslog"
	"os"
)

var (
	defaultConfig *logrusConfig
)

func logrusFormatter(f l.Format) logrus.Formatter {
	switch f {
	case l.JSON:
		return new(logrus.JSONFormatter)
	default:
		return &logrus.TextFormatter{ForceColors: false, DisableColors: true, FullTimestamp: true}
	}
}

func newField(key string, valType l.FieldType, val interface{}) logrusField {
	return logrusField{key: key, val: val, valType: valType}
}

type logrusField struct {
	key     string
	val     interface{}
	valType l.FieldType
}

func (f logrusField) Key() string {
	return f.key
}

func (f logrusField) Val() interface{} {
	return f.val
}

func (f logrusField) Type() l.FieldType {
	return f.valType
}

func newFieldProvider() *logrusFieldProvider {
	return new(logrusFieldProvider)
}

type logrusFieldProvider struct {
}

func (logrusFieldProvider) String(key string, val string) l.Field {
	return newField(key, l.StringField, val)
}

func (logrusFieldProvider) Bytes(key string, val []byte) l.Field {
	return newField(key, l.BytesField, val)
}

func (logrusFieldProvider) Int(key string, val int) l.Field {
	return newField(key, l.IntField, val)
}

func (logrusFieldProvider) Int32(key string, val int32) l.Field {
	return newField(key, l.Int32Field, val)
}

func (logrusFieldProvider) Int64(key string, val int64) l.Field {
	return newField(key, l.Int64Field, val)
}

func (logrusFieldProvider) Float(key string, val float32) l.Field {
	return newField(key, l.FloatField, val)
}

func (logrusFieldProvider) Float64(key string, val float64) l.Field {
	return newField(key, l.Float64Field, val)
}

func (logrusFieldProvider) Duration(key string, val time.Duration) l.Field {
	return newField(key, l.DurationField, val)
}

func (logrusFieldProvider) Time(key string, val time.Time) l.Field {
	return newField(key, l.TimeField, val)
}

func (logrusFieldProvider) Bool(key string, val bool) l.Field {
	return newField(key, l.BoolField, val)
}

func (logrusFieldProvider) Struct(key string, val interface{}) l.Field {
	return newField(key, l.StructField, val)
}

func (logrusFieldProvider) Slice(key string, val interface{}) l.Field {
	return newField(key, l.SliceField, val)
}

func (logrusFieldProvider) Error(val error) l.Field {
	return newField("error", l.ErrorField, val)
}

type logrusLogger struct {
	l.BaseLogger
	logger *logrus.Entry
}

func (l *logrusLogger) WithFields(fields ...l.Field) l.Logger {
	return &logrusLogger{
		logger: l.logger.WithFields(toLogrusFields(fields...)),
	}
}

func (l logrusLogger) Debug(message string, fields ...l.Field) {
	// if l.logger.Level < logrus.DebugLevel {
	// 	return
	// }
	if len(fields) <= 0 {
		l.logger.Debug(message)
	} else {
		l.logger.WithFields(toLogrusFields(fields...)).Debug(message)
	}
}

func (l logrusLogger) Info(message string, fields ...l.Field) {
	// if l.logger.Level < logrus.InfoLevel {
	// 	return
	// }
	if len(fields) <= 0 {
		l.logger.Info(message)
	} else {
		l.logger.WithFields(toLogrusFields(fields...)).Info(message)
	}
}

func (l logrusLogger) Warn(message string, fields ...l.Field) {
	// if l.logger.Level < logrus.WarnLevel {
	// 	return
	// }
	if len(fields) <= 0 {
		l.logger.Warn(message)
	} else {
		l.logger.WithFields(toLogrusFields(fields...)).Warn(message)
	}
}

func (l logrusLogger) Error(message string, fields ...l.Field) {
	if len(fields) <= 0 {
		l.logger.Error(message)
	} else {
		l.logger.WithFields(toLogrusFields(fields...)).Error(message)
	}
}

func (l logrusLogger) Panic(message string, fields ...l.Field) {
	if len(fields) <= 0 {
		l.logger.Panic(message)
	} else {
		l.logger.WithFields(toLogrusFields(fields...)).Panic(message)
	}
}

func (l logrusLogger) Fatal(message string, fields ...l.Field) {
	if len(fields) <= 0 {
		l.logger.Fatal(message)
	} else {
		l.logger.WithFields(toLogrusFields(fields...)).Fatal(message)
	}
}

func (logrusLogger) String() string {
	return "provider=logrus"
}

func Setup(loggerConfig *l.Configuration) error {
	logrusConfig, errs := toLogrusConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		if loggerConfig.Debug {
			fmt.Printf("l.SetupLogrusConfigErr Config=%s Errs=%v\n", logrusConfig.String(), errs)
		}
		return fmt.Errorf("l.SetupLogrusErr Errs=%v", errs)
	}
	logrus.SetLevel(logrusConfig.level)
	logrus.SetFormatter(logrusConfig.formatter)
	logrus.SetOutput(logrusConfig.output)
	defaultConfig = logrusConfig
	return l.Setup(New, newFieldProvider())
}

func New(field ...l.Field) (l.Logger, error) {
	return create(defaultConfig, field...)
}

func NewByConfig(loggerConfig *l.Configuration, field ...l.Field) (l.Logger, error) {
	logrusConfig, errs := toLogrusConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		if loggerConfig.Debug {
			fmt.Printf("l.logrus.NewByConfigErr Config=%s errs=%s\n", logrusConfig.String(), errs)
		}
		return nil, fmt.Errorf("l.logrus.NewByConfigErr Config=%s Errs=%v", logrusConfig.String(), errs)
	}
	return create(logrusConfig, field...)
}

func toLogrusFields(fields ...l.Field) logrus.Fields {
	logrusFields := make(map[string]interface{})
	for _, v := range fields {
		logrusFields[v.Key()] = v.Val()
	}
	return logrusFields
}

func create(cfg *logrusConfig, field ...l.Field) (l.Logger, error) {
	if cfg.debug {
		fmt.Printf("l.logrus.Creating Config=%s\n", cfg.String())
	}
	logBackend := &logrus.Logger{
		Level:     cfg.level,
		Formatter: cfg.formatter,
		Hooks:     make(logrus.LevelHooks),
		Out:       cfg.output,
	}
	// for _, hook := range cfg.hooks {
	// 	logBackend.Hooks.Add(hook)
	// }
	logger := new(logrusLogger)
	logger.logger = logrus.NewEntry(logBackend)
	logger.logger.Data = toLogrusFields(field...)
	if cfg.debug {
		fmt.Printf("l.logrus.Created Config=%s Logger=%s\n", cfg.String(), logger.String())
	}
	return logger, nil
}

type logrusConfig struct {
	debug     bool
	output    io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func (l logrusConfig) String() string {
	return fmt.Sprintf("debug=%t level=%s hasFormatter=%t hasOutput=%t", l.debug, l.level.String(), l.formatter != nil, l.output != nil)
}

func toLogrusConfig(cfg *l.Configuration) (*logrusConfig, []error) {
	var errs []error
	var output io.Writer
	if cfg.Out == l.Out("") {
		output = os.Stdout
	} else if tmpWriter, tmpErr := cfg.Out.GetOutput(); tmpErr != nil {
		errs = append(errs, tmpErr)
		output = os.Stdout
	} else {
		output = tmpWriter
	}
	var level logrus.Level
	if cfg.Level == l.Level("") {
		level = logrus.DebugLevel
	} else if tmlLevel, tmpErr := logrus.ParseLevel(cfg.Level.String()); tmpErr != nil {
		errs = append(errs, tmpErr)
		level = logrus.DebugLevel
	} else {
		level = tmlLevel
	}
	return &logrusConfig{
		debug:     cfg.Debug,
		level:     level,
		formatter: logrusFormatter(cfg.Format),
		output:    output,
	}, errs
}
