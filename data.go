package logger

import (
	"errors"
	"fmt"
	"strings"
)

const (
	//STDOUT any message to stdout
	STDOUT Out = "stdout"
	//STDERR redirects any message to stderr
	STDERR Out = "stderr"
	//DISCARD set logger to ignore all message
	DISCARD Out = "discard"

	//LOGRUS is the github.com/Sirupsen/logrus id
	LOGRUS Provider = "logrus"
	//ZAP is the github.com/uber-go/zap id
	ZAP Provider = "zap"

	//TEXT is the text log format
	TEXT Format = "text"
	//TEXTColor is the text log format with color
	TEXTColor Format = "text_color"
	//JSON is the json log format
	JSON Format = "json"
	//JSONColor is the json log format with color
	JSONColor Format = "json_color"
	//LOGRUSFmtfText is the text with the logrus formatf approach
	LOGRUSFmtfText Format = "logrusFrmtfText"

	//PANIC is the panic level logger
	PANIC Level = "panic"
	//FATAL is the fatal level logger
	FATAL Level = "fatal"
	//ERROR is the error level logger
	ERROR Level = "error"
	//WARN is the warn level logger
	WARN Level = "warn"
	//INFO is the info level logger
	INFO Level = "info"
	//DEBUG is the debug level logger
	DEBUG Level = "debug"

	//StringField is a constant for string logger fields
	StringField FieldType = iota
	//BytesField is a constant for byte slice logger fields
	BytesField
	//IntField is a constant for string logger fields
	IntField
	//Int64Field is a constant for string logger fields
	Int64Field
	//FloatField is a constant for string logger fields
	FloatField
	//Float64Field is a constant for string logger fields
	Float64Field
	//DurationField is a constant for duration logger fields
	DurationField
	//TimeField is a constant for time logger fields
	TimeField
	//BoolField is a constant for string logger fields
	BoolField
	//StructField is a constant for string logger fields
	StructField
	//ErrorField is a constant for error logger fields
	ErrorField
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was select
	ErrInvalidProvider = errors.New("Logger.InvalidProvider[Message='Avaible providers are: LOGRUS and ZAP']")
	//DefaultConfig holds the instance of the behavior parameters
	DefaultConfig *Configuration
)

//Provider is the back end implementor id of the logging feature
type Provider string

func (p Provider) String() string {
	return string(p)
}

// Set is a utility method for flag system usage
func (p *Provider) Set(value string) error {
	*p = Provider(value)
	return nil
}

//Out is the type for logger writer config
type Out string

func (o Out) String() string {
	return string(o)
}

// Set is a utility method for flag system usage
func (o *Out) Set(value string) error {
	*o = Out(value)
	return nil
}

//Hooks is the type to configure an create hooks for the logger implementation
type Hooks string

func (h Hooks) String() string {
	return string(h)
}

// Set is a utility method for flag system usage
func (h *Hooks) Set(value string) error {
	*h = Hooks(strings.TrimSpace(value))
	return nil
}

// Option is used to set options for the logger.
type Option interface {
	apply(Logger)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(Logger) error

func (f optionFunc) apply(l Logger) {
	f(l)
}

//Level is the threshold of the logger
type Level string

// String returns a lower-case ASCII representation of the log level.
func (n Level) String() string {
	return string(n)
}

// Set is a utility method for flag system usage
func (l *Level) Set(value string) error {
	*l = Level(value)
	return nil
}

//Format is a parameter to controle the logger style
type Format string

func (f Format) String() string {
	return string(f)
}

// Set is a utility method for flag system usage
func (f *Format) Set(value string) error {
	*f = Format(value)
	return nil
}

//Configuration holds the log beahvior parameters
type Configuration struct {
	Debug    bool     `json:"debug" mapstructure:"debug"`
	Provider Provider `json:"provider" mapstructure:"provider"`
	Level    Level    `json:"level" mapstructure:"level"`
	Format   Format   `json:"format" mapstructure:"format"`
	Out      Out      `json:"out" mapstructure:"out"`
	Hooks    Hooks    `json:"hooks" mapstructure:"hooks"`
}

func (l Configuration) String() string {
	return fmt.Sprintf("Configuration Provider=%s Level=%s Format=%s Out=%s Hooks=%s", l.Provider, l.Level, l.Format, l.Out, l.Hooks)
}

//FieldType is a type identifier for logger fields
type FieldType int8

//Field is a struct to send paramaters to log messages
type Field struct {
	key     string
	val     interface{}
	valType FieldType
}

//Logger is an interface to write log messages
type Logger interface {
	Level() Level
	IsEnabled(Level) bool
	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
	Fatal(string, ...Field)

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Panicf(string, ...interface{})
	Fatalf(string, ...interface{})
}

type baseLogger struct {
	level Level
}

func (b baseLogger) Level() Level {
	return b.level
}

func (b baseLogger) IsEnabled(level Level) bool {
	return b.level >= level
}

func (b baseLogger) Debugf(string, ...interface{}) {

}

func (b baseLogger) Infof(string, ...interface{}) {

}

func (b baseLogger) Warnf(string, ...interface{}) {

}

func (b baseLogger) Errorf(string, ...interface{}) {

}

func (b baseLogger) Panicf(string, ...interface{}) {

}

func (b baseLogger) Fatalf(string, ...interface{}) {

}
