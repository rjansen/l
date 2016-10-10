package logger

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"io"
	"log/syslog"
	"os"
	//"strings"
)

var (
	rootLogrusConfig *logrusConfig
)

func (f Format) toLogrusFormatter() logrus.Formatter {
	switch f {
	case JSON:
		return new(logrus.JSONFormatter)
	default:
		return &logrus.TextFormatter{ForceColors: false, DisableColors: true, FullTimestamp: true}
	}
}

func (o Out) toLogrusOut() (io.Writer, error) {
	return getOutput(o)
}

func (n Level) toLogrusLevel() (logrus.Level, error) {
	return logrus.ParseLevel(n.String())
}

func (h Hooks) toLogrusSyslogHook() (*logrus_syslog.SyslogHook, error) {
	hookName := string(h)
	switch hookName {
	case "syslog":
		//hook, err := logrus_syslog.NewSyslogHook("udp", "127.0.0.1:514", syslog.LOG_DEBUG, "glive")
		hook, err := logrus_syslog.NewSyslogHook("udp", "127.0.0.1:514", syslog.LOG_DEBUG, "glive")
		if err != nil {
			return nil, err
		}
		return hook, nil
	default:
		return nil, nil
	}

	// hooksValue := string(h)
	// if hooksValue == "" {
	// 	return nil, nil
	// }
	// var hookName string
	// //var hookParams string
	// if !strings.Contains(hooksValue, "?") {
	// 	//hookName, hookParams = hooksValue, ""
	// 	hookName = hooksValue
	// } else {
	// 	hookConfig := strings.Split(hooksValue, "?")
	// 	//hookName, hookParameters := hookConfig[0], hookConfig[1]
	// 	//hookName, hookParams = hookConfig[0], hookConfig[1]
	// 	hookName = hookConfig[0]
	// }
	// var hooks []logrus.Hook
	// switch hookName {
	// case "syslog":
	// 	//hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", syslog.LOG_DEBUG, "glive_localhost")
	// 	hook, err := logrus_syslog.NewSyslogHook("tcp", "127.0.0.1:514", syslog.LOG_DEBUG, "glive_localhost")
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	hooks = append(hooks, hook)
	// 	return hooks, nil
	// default:
	// 	return nil, nil
	// }
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
	if l.logger.Level < logrus.DebugLevel {
		return
	}
	if len(fields) <= 0 {
		l.logger.Debug(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Debug(message)
	}
}

func (l logrusLogger) Info(message string, fields ...Field) {
	if l.logger.Level < logrus.InfoLevel {
		return
	}
	if len(fields) <= 0 {
		l.logger.Info(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Info(message)
	}
}

func (l logrusLogger) Warn(message string, fields ...Field) {
	if l.logger.Level < logrus.WarnLevel {
		return
	}
	if len(fields) <= 0 {
		l.logger.Warn(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Warn(message)
	}
}

func (l logrusLogger) Error(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Error(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Error(message)
	}
}

func (l logrusLogger) Panic(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Panic(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Panic(message)
	}
}

func (l logrusLogger) Fatal(message string, fields ...Field) {
	if len(fields) <= 0 {
		l.logger.Fatal(message)
	} else {
		l.logger.WithFields(l.toLogrusFields(fields...)).Fatal(message)
	}
}

func (l logrusLogger) Debugf(message string, fields ...interface{}) {
	l.logger.Debugf(message, fields...)
}

func (l logrusLogger) Infof(message string, fields ...interface{}) {
	l.logger.Infof(message, fields...)
}

func (l logrusLogger) Warnf(message string, fields ...interface{}) {
	l.logger.Warnf(message, fields...)
}

func (l logrusLogger) Errorf(message string, fields ...interface{}) {
	l.logger.Errorf(message, fields...)
}

func (l logrusLogger) Panicf(message string, fields ...interface{}) {
	l.logger.Panicf(message, fields...)
}

func (l logrusLogger) Fatalf(message string, fields ...interface{}) {
	l.logger.Fatalf(message, fields...)
}

func setupLogrus(loggerConfig Configuration) error {
	logrusConfig, errs := createLogrusConfig(loggerConfig)
	if errs != nil && len(errs) > 0 {
		return fmt.Errorf("SetupLogrusErr[Errs=%v]", errs)
	}
	rootLogrusConfig = &logrusConfig
	logrus.SetLevel(rootLogrusConfig.level)
	logrus.SetFormatter(rootLogrusConfig.formatter)
	logrus.SetOutput(rootLogrusConfig.output)
	loggerFactory = newLogrus
	return nil
}

func newLogrus(config Configuration) Logger {
	logrusConfig, errs := createLogrusConfig(config)
	if errs != nil {
		fmt.Printf("NewLogrusConfigErr=%+v\n", errs)
	}
	if config.Debug {
		fmt.Printf("NewLogrusConfig=%s\n", logrusConfig.String())
	}
	logger := new(logrusLogger)
	logger.logger = &logrus.Logger{
		Level:     logrusConfig.level,
		Formatter: logrusConfig.formatter,
		Hooks:     make(logrus.LevelHooks),
		Out:       logrusConfig.output,
	}
	//logger.logger.SetNoLock()
	//l3, err := syslog.Dial("udp", "127.0.0.1:514", syslog.LOG_ERR, "glive")
	//l3, err := syslog.Dial("udp", "localhost", syslog.LOG_ERR, "GoExample") // connection to a log daemon
	//defer l3.Close()
	//if err != nil {
	//	fmt.Println("CreateSyslogErr", err.Error())
	//}
	//l3.Err("SyslogSimpleMsgErr")
	hooks, err := config.Hooks.toLogrusSyslogHook()
	if err != nil {
		fmt.Println("CreateSyslogHook", err.Error())
	}
	if hooks != nil {
		//for _, hook := range hooks {
		hooks.Writer.Err("SettingSyslogErr")
		logger.logger.Hooks.Add(hooks)
		fmt.Println("SyslogHookAdded")
		//}
	}
	if config.Debug {
		fmt.Printf("NewLogrusLogger=%+v\n", logger.logger)
	}
	return logger
}

type logrusConfig struct {
	output    io.Writer
	formatter logrus.Formatter
	level     logrus.Level
}

func (l logrusConfig) String() string {
	return fmt.Sprintf("logrusConfig[level=%s formatter=%t output=%t]", l.level.String(), l.formatter != nil, l.output != nil)
}

func createLogrusConfig(cfg Configuration) (logrusConfig, []error) {
	var errs []error
	var output io.Writer
	if cfg.Out == Out("") {
		output = os.Stdout
	} else if tmpWriter, tmpErr := cfg.Out.toLogrusOut(); tmpErr != nil {
		errs = append(errs, tmpErr)
		output = os.Stdout
	} else {
		output = tmpWriter
	}
	var level logrus.Level
	if cfg.Level == Level("") {
		level = logrus.DebugLevel
	} else if tmlLevel, tmpErr := logrus.ParseLevel(cfg.Level.String()); tmpErr != nil {
		errs = append(errs, tmpErr)
		level = logrus.DebugLevel
	} else {
		level = tmlLevel
	}
	return logrusConfig{
		level:     level,
		formatter: cfg.Format.toLogrusFormatter(),
		output:    output,
	}, errs
}
