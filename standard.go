package l

import (
	"io"
	"log"
	"os"
	"time"
)

func newField(key string, valType FieldType, val interface{}) field {
	return field{key: key, val: val, valType: valType}
}

type field struct {
	key     string
	val     interface{}
	valType FieldType
}

func (f field) Key() string {
	return f.key
}

func (f field) Val() interface{} {
	return f.val
}

func (f field) Type() FieldType {
	return f.valType
}

func newFieldProvider() *stdFieldProvider {
	return new(stdFieldProvider)
}

type stdFieldProvider struct {
}

func (stdFieldProvider) String(key string, val string) Field {
	return newField(key, StringField, val)
}

func (stdFieldProvider) Bytes(key string, val []byte) Field {
	return newField(key, BytesField, val)
}

func (stdFieldProvider) Int(key string, val int) Field {
	return newField(key, IntField, val)
}

func (stdFieldProvider) Int32(key string, val int32) Field {
	return newField(key, Int32Field, val)
}

func (stdFieldProvider) Int64(key string, val int64) Field {
	return newField(key, Int64Field, val)
}

func (stdFieldProvider) Float(key string, val float32) Field {
	return newField(key, FloatField, val)
}

func (stdFieldProvider) Float64(key string, val float64) Field {
	return newField(key, Float64Field, val)
}

func (stdFieldProvider) Duration(key string, val time.Duration) Field {
	return newField(key, DurationField, val)
}

func (stdFieldProvider) Time(key string, val time.Time) Field {
	return newField(key, TimeField, val)
}

func (stdFieldProvider) Bool(key string, val bool) Field {
	return newField(key, BoolField, val)
}

func (stdFieldProvider) Struct(key string, val interface{}) Field {
	return newField(key, StructField, val)
}

func (stdFieldProvider) Slice(key string, val interface{}) Field {
	return newField(key, SliceField, val)
}

func (stdFieldProvider) Error(val error) Field {
	return newField("error", ErrorField, val)
}

//standardLogger holds the level loggers pointer
type standardLogger struct {
	BaseLogger
	d *log.Logger
	i *log.Logger
	w *log.Logger
	e *log.Logger
	p *log.Logger
	f *log.Logger
}

func (s *standardLogger) WithFields(fields ...Field) Logger {
	return s
}

func (s *standardLogger) Debug(m string, fields ...Field) {
	s.d.Printf("message=%s fields=%+v", m, fields)
}

func (s *standardLogger) Info(m string, fields ...Field) {
	s.i.Printf("message=%s fields=%+v", m, fields)
}

func (s *standardLogger) Warn(m string, fields ...Field) {
	s.w.Printf("message=%s fields=%+v", m, fields)
}

func (s *standardLogger) Error(m string, fields ...Field) {
	s.e.Printf("message=%s fields=%+v", m, fields)
}

func (s *standardLogger) Panic(m string, fields ...Field) {
	s.f.Printf("message=%s fields=%+v", m, fields)
}

func (s *standardLogger) Fatal(m string, fields ...Field) {
	s.f.Printf("message=%s fields=%+v", m, fields)
}

func (standardLogger) String() string {
	return "provider=standard"
}

func newLogger(field ...Field) (Logger, error) {
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
