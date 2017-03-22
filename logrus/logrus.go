package logrus

import (
	"fmt"
	"github.com/Sirupsen/logrus"
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

type logrusLogger struct {
	l.BaseLogger
	logger *logrus.Logger
}

func (l logrusLogger) toLogrusFields(fields ...l.Field) logrus.Fields {
	logrusFields := make(map[string]interface{})
	for _, v := range fields {
		logrusFields[v.Key] = v.Val
	}
	return logrusFields
}

func (l logrusLogger) Debug(message string, fields ...l.Field) {
	if l.logger.Level < logrus.DebugLevel {
		return
	}
	if len(fields) <= 0 {
		l.logger.Debug(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Debug(message)
	}
}

func (l logrusLogger) Info(message string, fields ...l.Field) {
	if l.logger.Level < logrus.InfoLevel {
		return
	}
	if len(fields) <= 0 {
		l.logger.Info(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Info(message)
	}
}

func (l logrusLogger) Warn(message string, fields ...l.Field) {
	if l.logger.Level < logrus.WarnLevel {
		return
	}
	if len(fields) <= 0 {
		l.logger.Warn(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Warn(message)
	}
}

func (l logrusLogger) Error(message string, fields ...l.Field) {
	if len(fields) <= 0 {
		l.logger.Error(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Error(message)
	}
}

func (l logrusLogger) Fatal(message string, fields ...l.Field) {
	if len(fields) <= 0 {
		l.logger.Fatal(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Fatal(message)
	}
}

func Setup(loggerConfig *Configuration) error {
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
	return l.Setup(New)
}

func New(field ...l.Field) l.Logger {
	if defaultConfig.debug {
		fmt.Printf("l.CreatingLogrusLogger Config=%s\n", defaultConfig.String())
	}
	logger := create(defaultConfig)
	if defaultConfig.debug {
		fmt.Printf("l.LogrusLoggerCreated Logger=%+v\n", logger)
	}
	return logger
}

func NewByConfig(loggerConfig *Configuration, field ...l.Field) (l.Logger, error) {
	logrusConfig, errs := toLogrusConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		if loggerConfig.Debug {
			fmt.Printf("l.SetupLogrusConfigErr Config=%s errs=%s\n", logrusConfig.String(), errs)
		}
		return nil, fmt.Errorf("l.NewLogrusByConfigErr Errs=%v", errs)
	}
	if loggerConfig.Debug {
		fmt.Printf("l.CreatingLogrusByConfigLogger Config=%s\n", logrusConfig.String())
	}
	logger := new(logrusLogger)
	logger.logger = &logrus.Logger{
		Level:     logrusConfig.level,
		Formatter: logrusConfig.formatter,
		Hooks:     make(logrus.LevelHooks),
		Out:       logrusConfig.output,
	}
	if defaultConfig.debug {
		fmt.Printf("l.LogrusLoggerCreatedByConfig Logger=%+v\n", logger.logger)
	}
	return logger, nil
}

func create(cfg *logrusConfig) l.Logger {
	logger := new(logrusLogger)
	logger.logger = &logrus.Logger{
		Level:     cfg.level,
		Formatter: cfg.formatter,
		Hooks:     make(logrus.LevelHooks),
		Out:       cfg.output,
	}
	return logger

}

type logrusConfig struct {
	debug     bool
	output    io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func (l logrusConfig) String() string {
	return fmt.Sprintf("logrusConfig level=%s formatter=%t output=%t", l.level.String(), l.formatter != nil, l.output != nil)
}

func toLogrusConfig(cfg *Configuration) (*logrusConfig, []error) {
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
