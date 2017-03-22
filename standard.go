package l

import (
	"io"
	"log"
	"os"
)

//standardLogger holds the level loggers pointer
type standardLogger struct {
	BaseLogger
	d *log.Logger
	i *log.Logger
	w *log.Logger
	e *log.Logger
	f *log.Logger
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

func (s *standardLogger) Fatal(m string, fields ...Field) {
	s.f.Printf("message=%s fields=%+v", m, fields)
}

func newLogger(field ...Field) Logger {
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
		f: log.New(errorWriter, "FATAL ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return _logger
}
