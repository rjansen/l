package l

import (
	"io"
	"log"
	"os"
)

//DefaultLogger holds the level loggers pointer
type DefaultLogger struct {
	Debug *log.Logger
	Info  *log.Logger
	Note  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Panic *log.Logger
}

func newLogger(loggerConfig *Configuration) *DefaultLogger {
	//fmt.Printf("CreatingLogger: File=%v Level=%v\n", loggerConfig.Output, loggerConfig.Level)
	output, err := os.OpenFile(string(loggerConfig.Out), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		//fmt.Printf("CreateOrOpenLoggerFileError: Message='%v'", err)
	}
	errorWriter := io.MultiWriter(output, os.Stderr)
	_logger := &DefaultLogger{
		Debug: log.New(output, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(output, "INFO ", log.Ldate|log.Ltime|log.Lshortfile),
		Note:  log.New(output, "NOTE ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(output, "WARN ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(errorWriter, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
		Panic: log.New(errorWriter, "PANIC ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return _logger
}
