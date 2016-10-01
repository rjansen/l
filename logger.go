package logger

import (
	"errors"
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	logging "github.com/op/go-logging"
	zap "github.com/uber-go/zap"
	"io"
	"log"
	"os"
)

const (
	//LOGRUS is the github.com/Sirupsen/logrus id
	LOGRUS Provider = "logrus"
	//ZAP is the github.com/uber-go/zap id
	ZAP = "zap"
	//OP is the github.com/op/go-logging id
	OP = "op"

	//PANIC is the panic level logger
	PANIC Level = iota
	//ERROR is the error level logger
	ERROR
	//WARN is the warn level logger
	WARN
	//INFO is the info level logger
	INFO
	//DEBUG is the debug level logger
	DEBUG
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was select
	ErrInvalidProvider = errors.New("Logger.InvalidProvider[Message='Avaible providers are: LOGRUS, ZAP and OP']")
	//Config holds the instance of the behavior parameters
	Config     *Configuration
	loggerFile *os.File
)

//Level is the threshold of the logger
type Level int

//Int cast the Level into an int representation
//func (l Level) Int() int {
//	return int(l)
//}

//Provider is the back end implementor id of the logging feature
type Provider string

//Configuration holds the log beahvior parameters
type Configuration struct {
	File         string
	DefaultLevel Level
	Provider     Provider
	Format       string
}

func (l *Configuration) String() string {
	return fmt.Sprintf("Config[File=%v DefaultLevel=%v Format=%v]", l.File, l.DefaultLevel, l.Format)
}

//Setup initializes the logger system
func Setup(loggerConfig *Configuration) error {
	switch loggerConfig.Provider {
	case LOGRUS:
		return setupLogrus(loggerConfig)
	case ZAP:
		return setupZap(loggerConfig)
	case OP:
		return setupOp(loggerConfig)
	default:
		return ErrInvalidProvider

	}
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

func setupLogrus(loggerConfig *Configuration) error {
	Config = loggerConfig
	return nil
}

func setupZap(loggerConfig *Configuration) error {
	Config = loggerConfig
	return nil
}

func setupOp(loggerConfig *Configuration) error {
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
	if loggerConfig.DefaultLevel < 0 || loggerConfig.DefaultLevel > 5 {
		loggerLevel = logging.DEBUG
	} else {
		loggerLevel = logging.Level(loggerConfig.DefaultLevel)
	}
	levelMessages.SetLevel(loggerLevel, "")

	levelError := logging.AddModuleLevel(backEndError)
	levelError.SetLevel(logging.ERROR, "")

	logging.SetBackend(levelMessages, levelError)
	//fmt.Printf("LoggerConfiguredSuccessfully: LoggerConfig=%v\n", loggerConfig)
	Config = loggerConfig
	return nil
}

func newLogrus(loggerConfig *Configuration) *logrus.Logger {
	return logrus.New()
}

func newZap(loggerConfig *Configuration) zap.Logger {
	return zap.New(zap.NewTextEncoder())
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

func newLogger(loggerConfig *Configuration) *Logger {
	fmt.Printf("CreatingLogger: File=%v DefaultLevel=%v\n", loggerConfig.File, loggerConfig.DefaultLevel)
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

func getLoggerFile(loggerConfig *Configuration) (*os.File, error) {
	fmt.Printf("CreatingLoggerFile: File=%v DefaultLevel=%v\n", loggerConfig.File, loggerConfig.DefaultLevel)
	loggerFile, err := os.OpenFile(loggerConfig.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("CreateOrOpenLoggerFileError: Message='%v'\n", err)
		return nil, fmt.Errorf("CreateOrOpenLoggerFileError: Message='%v'", err)
	}
	return loggerFile, nil
}
