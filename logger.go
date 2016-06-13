package logger

import (
	"farm.e-pedion.com/repo/config"
	"fmt"
	logging "github.com/op/go-logging"
	"io"
	"log"
	"os"
)

var (
	loggerFile *os.File
)

//Setup initializes the logger system
func Setup(loggerConfig *config.LoggerConfig) {
	var loggerWriter io.Writer
	loggerFile, err := getLoggerFile(loggerConfig)
	if err != nil {
		loggerWriter = os.Stdout
	} else {
		loggerWriter = loggerFile
	}
	loggerFormat := logging.MustStringFormatter(loggerConfig.Format)
	//TODO: Remove os.Stdout. For performance reasons the log messages must send only to the file
	backEndMessages := logging.NewBackendFormatter(logging.NewLogBackend(io.MultiWriter(os.Stdout, loggerWriter), "", 0), loggerFormat)
	backEndError := logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), loggerFormat)

	levelMessages := logging.AddModuleLevel(backEndMessages)
	var loggerLevel logging.Level
	if loggerConfig.Level < 0 || loggerConfig.Level > 5 {
		loggerLevel = logging.DEBUG
	} else {
		loggerLevel = logging.Level(loggerConfig.Level)
	}
	levelMessages.SetLevel(loggerLevel, "")

	levelError := logging.AddModuleLevel(backEndError)
	levelError.SetLevel(logging.ERROR, "")

	logging.SetBackend(levelMessages, levelError)
	fmt.Printf("LoggerConfiguredSuccessfully: LoggerConfig=%v\n", loggerConfig)
}

//Close closes the log file
func Close() {
	if loggerFile != nil {
		fmt.Printf("ClosingLoggerFile: File=%v\n", loggerFile.Name)
		loggerFile.Close()
	}
}

//GetLogger gets the go-logging underlying instance of the model provided
func GetLogger(module string) *logging.Logger {
	return logging.MustGetLogger(module)
}

func getLoggerFile(loggerConfig *config.LoggerConfig) (*os.File, error) {
	fmt.Printf("CreatingLoggerFile: File=%v Level=%v\n", loggerConfig.File, loggerConfig.Level)
	loggerFile, err := os.OpenFile(loggerConfig.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("CreateOrOpenLoggerFileError: Message='%v'\n", err)
		return nil, fmt.Errorf("CreateOrOpenLoggerFileError: Message='%v'", err)
	}
	return loggerFile, nil
}

//Logger holds the level loggers pointer
type Logger struct {
	Debug *log.Logger
	Info  *log.Logger
	Note  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Panic *log.Logger
}

func newLogger(loggerConfig *config.LoggerConfig) *Logger {
	fmt.Printf("CreatingLogger: File=%v Level=%v\n", loggerConfig.File, loggerConfig.Level)
	loggerFile, err := os.OpenFile(loggerConfig.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("CreateOrOpenLoggerFileError: Message='%v'", err)
	}
	errorWriter := io.MultiWriter(loggerFile, os.Stderr)
	_logger := &Logger{
		Debug: log.New(loggerFile, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile),
		Info:  log.New(loggerFile, "INFO ", log.Ldate|log.Ltime|log.Lshortfile),
		Note:  log.New(loggerFile, "NOTE ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:  log.New(loggerFile, "WARN ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(errorWriter, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile),
		Panic: log.New(errorWriter, "PANIC ", log.Ldate|log.Ltime|log.Lshortfile),
	}
	return _logger
}
