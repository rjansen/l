package l

import "context"

const (
	//STDOUT any message to stdout
	STDOUT Out = "stdout"
	//STDERR redirects any message to stderr
	STDERR Out = "stderr"

	//ERROR is the error level logger
	ERROR Level = "error"
	//INFO is the info level logger
	INFO Level = "info"
	//DEBUG is the debug level logger
	DEBUG Level = "debug"
)

//Out is the type for logger writer config
type Out string

func (o Out) String() string {
	return string(o)
}

// Set is a utility method for flag system usage
func (o *Out) Set(value string) error {
	switch value {
	case "stdout", "STDOUT", "":
		*o = STDOUT
	case "stderr", "STDERR":
		*o = STDERR
	default:
		*o = Out(value)
	}
	return nil
}

//Level is the threshold of the logger
type Level string

// String returns a lower-case ASCII representation of the log level.
func (l Level) String() string {
	return string(l)
}

// Set is a utility method for flag system usage
func (l *Level) Set(value string) error {
	switch value {
	case "info", "INFO":
		*l = INFO
	case "error", "ERROR":
		*l = ERROR
	default:
		*l = DEBUG
	}
	return nil
}

type Value struct {
	name  string
	value interface{}
}

func NewValue(name string, value interface{}) Value {
	return Value{name: name, value: value}
}

type Logger interface {
	Debug(context.Context, string, ...Value)
	Info(context.Context, string, ...Value)
	Error(context.Context, string, ...Value)
}

type LogWriter interface {
	Write(...Value)
}

type Driver interface {
	Log(Level, string) LogWriter
	Close()
}

type logger struct {
	driver Driver
}

func (log logger) log(_ context.Context, level Level, msg string, values ...Value) {
	if writer := log.driver.Log(level, msg); writer != nil {
		writer.Write(values...)
	}
}

func (log logger) Debug(ctx context.Context, msg string, values ...Value) {
	log.log(ctx, DEBUG, msg, values...)
}

func (log logger) Info(ctx context.Context, msg string, values ...Value) {
	log.log(ctx, INFO, msg, values...)
}

func (log logger) Error(ctx context.Context, msg string, values ...Value) {
	log.log(ctx, ERROR, msg, values...)
}

func New(driver Driver) Logger {
	return logger{
		driver: driver,
	}
}
