package l

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	//STDOUT any message to stdout
	STDOUT Out = "stdout"
	//STDERR redirects any message to stderr
	STDERR Out = "stderr"
	//DISCARD set logger to ignore all message
	DISCARD Out = "discard"

	//TEXT is the text log format
	TEXT Format = "text"
	//TEXTColor is the text log format with color
	TEXTColor Format = "text_color"
	//JSON is the json log format
	JSON Format = "json"
	//JSONColor is the json log format with color
	JSONColor Format = "json_color"

	//FATAL is the fatal level logger
	FATAL Level = "fatal"
	//PANIC is the panic level logger
	PANIC Level = "panic"
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
	//IntField is a constant for int logger fields
	IntField
	//Int32Field is a constant for int32 logger fields
	Int32Field
	//Int64Field is a constant for int64 logger fields
	Int64Field
	//FloatField is a constant for float32 logger fields
	FloatField
	//Float64Field is a constant for float64 logger fields
	Float64Field
	//DurationField is a constant for duration logger fields
	DurationField
	//TimeField is a constant for time logger fields
	TimeField
	//BoolField is a constant for bool logger fields
	BoolField
	//StructField is a constant for dynamic logger fields
	StructField
	//SliceField is a constant for slice logger fields
	SliceField
	//ErrorField is a constant for error logger fields
	ErrorField
)

//Hooks is the type to configure an create hooks for the logger implementation
type Hooks struct {
	Syslog SocketHook `json:"syslog" mapstructure:"syslog"`
	Gelf   SocketHook `json:"gelf" mapstructure:"gelf"`
	Stdout bool       `json:"stdout" mapstructure:"stdout"`
}

func (h Hooks) String() string {
	return fmt.Sprintf("Syslog=%s Gelf=%s Stdout=%t", h.Syslog.String(), h.Gelf.String(), h.Stdout)
}

//SocketHook is a hook that intent to sends data over network sockets
type SocketHook struct {
	Socket  string `json:"socket" mapstructure:"socket"`
	Address string `json:"addr" mapstructure:"addr"`
	Level   string `json:"level" mapstructure:"level"`
}

func (s SocketHook) String() string {
	return fmt.Sprintf("Socket=%s Address=%s Level=%s", s.Socket, s.Address, s.Level)
}

//Configuration holds the log beahvior parameters
type Configuration struct {
	Debug  bool   `json:"debug" mapstructure:"debug"`
	Level  Level  `json:"level" mapstructure:"level"`
	Format Format `json:"format" mapstructure:"format"`
	Out    Out    `json:"out" mapstructure:"out"`
	Hooks  Hooks  `json:"hooks" mapstructure:"hooks"`
}

func (l Configuration) String() string {
	return fmt.Sprintf("Level=%s Format=%s Out=%s Hooks=%s", l.Level, l.Format, l.Out, l.Hooks)
}

//Provider is the contract for logger factories
type Provider func(...Field) (Logger, error)

//FieldProvider is the contract for logger fields factories
type FieldProvider interface {
	String(string, string) Field
	Bytes(string, []byte) Field
	Int(string, int) Field
	Int32(string, int32) Field
	Int64(string, int64) Field
	Float(string, float32) Field
	Float64(string, float64) Field
	Duration(string, time.Duration) Field
	Time(string, time.Time) Field
	Bool(string, bool) Field
	Struct(string, interface{}) Field
	Slice(string, interface{}) Field
	Error(error) Field
}

//FieldType is a type identifier for logger fields
type FieldType int8

//Field is a struct to send paramaters to log messages
type Field interface {
	Key() string
	Val() interface{}
	Type() FieldType
}

//Out is the type for logger writer config
type Out string

func (o Out) String() string {
	return string(o)
}

// Set is a utility method for flag system usage
func (o *Out) Set(value string) error {
	if strings.TrimSpace(value) != "" {
		*o = Out(value)
	} else {
		*o = STDOUT
	}
	return nil
}

//GetOutput returns the out writer instance
func (o Out) GetOutput() (io.Writer, error) {
	switch o {
	case STDOUT:
		return os.Stdout, nil
	case STDERR:
		return os.Stderr, nil
	case DISCARD:
		return ioutil.Discard, nil
	default:
		file, err := os.OpenFile(o.String(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("l.CreateFileOutputErr Out=%v Message='%v'", o, err)
		}
		return file, nil
	}
}

//Level is the threshold of the logger
type Level string

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	return string(l)
}

// Set is a utility method for flag system usage
func (l *Level) Set(value string) error {
	if strings.TrimSpace(value) != "" {
		*l = Level(value)
	} else {
		*l = DEBUG
	}
	return nil
}

//Format is a parameter to controle the logger style
type Format string

func (f Format) String() string {
	return string(f)
}

// Set is a utility method for flag system usage
func (f *Format) Set(value string) error {
	if strings.TrimSpace(value) != "" {
		*f = Format(value)
	} else {
		*f = TEXT
	}
	return nil
}

//Logger is an interface to write log messages
type Logger interface {
	Level() Level
	Enabled(Level) bool
	WithFields(...Field) Logger

	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Panic(string, ...Field)
	Fatal(string, ...Field)

	String() string
}

type BaseLogger struct {
	level Level
}

func (b BaseLogger) Level() Level {
	return b.level
}

func (b BaseLogger) Enabled(level Level) bool {
	return b.level >= level
}
