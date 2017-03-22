package l

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
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
	//SliceField is a constant for slice logger fields
	SliceField
	//ErrorField is a constant for error logger fields
	ErrorField
)

//Provider is the type for create loggers
type Provider func(...Field) Logger

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

//FieldType is a type identifier for logger fields
type FieldType int8

//Field is a struct to send paramaters to log messages
type Field struct {
	Key     string
	Val     interface{}
	ValType FieldType
}

//Logger is an interface to write log messages
type Logger interface {
	Level() Level
	Enabled(Level) bool

	Debug(string, ...Field)
	Info(string, ...Field)
	Warn(string, ...Field)
	Error(string, ...Field)
	Fatal(string, ...Field)
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
