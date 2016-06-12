package logger

import (
	"farm.e-pedion.com/repo/config"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	loggerConfig *config.LoggerConfig
	logger       *Logger
)

//Logger is a type used to write application messages with level configuration
type Logger struct {
	file  *os.File
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
}

//Close release the logger instance and close all related resources
func (l *Logger) Close() {
	//TODO: Check if the nil attribution is a good pratice and it is necessary
	l.Debug = nil
	l.Info = nil
	l.Warn = nil
	l.Error = nil
	if l.file != nil {
		l.file.Close()
	}
}

//Init initializes the logger package
func Init(configuration *config.LoggerConfig) {
	if loggerConfig == nil {
		loggerConfig = configuration
	}
}

//GetLogger creates, only if necessary, the singleton instance of the logger component and returns
func GetLogger() *Logger {
	if logger == nil {
		logger = newLogger()
	}
	return logger
}

func newLogger() *Logger {
	fmt.Printf("CreatingLogger: File=%v Level=%v", loggerConfig.LoggerFile, loggerConfig.LoggerLevel)
	loggerFile, err := os.OpenFile(loggerConfig.LoggerFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("CreateOrOpenLoggerFileError: Message='%v'", err)
	}
	_logger := &Logger{
		file:  loggerFile,
		Debug: log.New(loggerFile, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(loggerFile, "INFO ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(loggerFile, "WARN ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(io.MultiWriter(loggerFile, os.Stderr), "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return _logger
}
