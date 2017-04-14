package l

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func newField(name string, value interface{}) field {
	return field{name: name, value: value}
}

type field struct {
	name  string
	value interface{}
}

func (f field) Name() string {
	return f.name
}

func (f field) Value() interface{} {
	return f.value
}

func (f field) String() string {
	return fmt.Sprintf("%s=%v", f.Name(), f.Value())
}

func DefaultFieldAdapter() FieldAdapter {
	return new(defaultFieldAdapter)
}

type defaultFieldAdapter struct {
}

func (defaultFieldAdapter) String(key string, val string) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Bytes(key string, val []byte) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Int(key string, val int) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Int32(key string, val int32) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Int64(key string, val int64) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Float(key string, val float32) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Float64(key string, val float64) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Duration(key string, val time.Duration) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Time(key string, val time.Time) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Bool(key string, val bool) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Struct(key string, val interface{}) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Slice(key string, val interface{}) Field {
	return newField(key, val)
}

func (defaultFieldAdapter) Error(val error) Field {
	return newField("error", val)
}

//standardLogger holds the level loggers pointer
type standardLogger struct {
	BaseLogger
	fields []Field
	d      *log.Logger
	i      *log.Logger
	w      *log.Logger
	e      *log.Logger
	p      *log.Logger
	f      *log.Logger
}

func (s *standardLogger) WithFields(fields ...Field) Logger {
	s.fields = append(s.fields, fields...)
	return s
}

func (s *standardLogger) Debug(m string, fields ...Field) {
	s.d.Printf("message=%s %v", m, fields)
}

func (s *standardLogger) Info(m string, fields ...Field) {
	s.i.Printf("message=%s %v", m, fields)
}

func (s *standardLogger) Warn(m string, fields ...Field) {
	s.w.Printf("message=%s %v", m, fields)
}

func (s *standardLogger) Error(m string, fields ...Field) {
	s.e.Printf("message=%s %v", m, fields)
}

func (s *standardLogger) Panic(m string, fields ...Field) {
	s.f.Printf("message=%s %v", m, fields)
}

func (s *standardLogger) Fatal(m string, fields ...Field) {
	s.f.Printf("message=%s %v", m, fields)
}

func (standardLogger) String() string {
	return "provider=standard"
}

func DefaultLog(c *Configuration, field ...Field) (Logger, error) {
	//fmt.Printf("CreatingLogger: File=%v Level=%v\n", loggerConfig.Output, loggerConfig.Level)
	// output, err := os.OpenFile(string(loggerConfig.Out), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err != nil {
	//fmt.Printf("CreateOrOpenLoggerFileError: Message='%v'", err)
	// }
	errorWriter := io.MultiWriter(os.Stdout, os.Stderr)
	_logger := &standardLogger{
		d: log.New(os.Stdout, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile),
		i: log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile),
		w: log.New(os.Stdout, "WARN ", log.Ldate|log.Ltime|log.Lshortfile),
		e: log.New(errorWriter, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
		p: log.New(errorWriter, "PANIC ", log.Ldate|log.Ltime|log.Lshortfile),
		f: log.New(errorWriter, "FATAL ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return _logger, nil
}
