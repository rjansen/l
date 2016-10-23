package logger

import (
	"errors"
	"farm.e-pedion.com/repo/config"
	"fmt"
	"github.com/matryer/resync"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	//ErrInvalidProvider is the err raised when an invalid provider was select
	ErrInvalidProvider = errors.New("logger.InvalidProvider Message='Avaible providers are: LOGRUS and ZAP'")
	//ErrSetupNeverCalled is raised when the Setup method does not call
	ErrSetupNeverCalled = errors.New("logger.SetupNeverCalledErr Message='You must call logger.Setup before execute this action'")
	once                resync.Once
	loggerFactory       func(Configuration) Logger
	defaultConfig       *Configuration
	rootLogger          Logger
)

func init() {
	fmt.Printf("logger.init\n")
	// if err := config.UnmarshalKey("logger.root", &defaultConfig); err != nil {
	// 	panic(err)
	// }
	// if defaultConfig.Debug {
	// 	fmt.Printf("logger.init configuration=%s", defaultConfig.String())
	// }
}

//Setup initializes the logger system
func Setup(loggerConfig Configuration) error {
	if loggerConfig.Debug {
		fmt.Printf("logger.Setup config=%+v\n", loggerConfig)
	}
	var setupErr error
	switch loggerConfig.Provider {
	case ZAP:
		setupErr = setupZap(loggerConfig)
	default:
		setupErr = setupLogrus(loggerConfig)
	}
	if setupErr != nil {
		return setupErr
	}
	return nil
}

//Get gets an implemetor of the configured log provider
func Get() Logger {
	once.Do(func() {
		if rootLogger == nil {
			setted := isSetted()
			if !setted {
				fmt.Printf("logger.Get.isSetted setted=%t defaultConfigIsNil=%t loggerFactoryIsNil=%t\n", setted, defaultConfig == nil, loggerFactory == nil)
				loggerConfig, _ := getConfiguration("logger.root")
				if err := Setup(loggerConfig); err != nil {
					panic(err)
				}
			}
			rootLogger = create()
		}
	})
	return rootLogger
}

func isSetted() bool {
	return loggerFactory != nil && defaultConfig != nil
}

func create() Logger {
	setted := isSetted()
	if !setted {
		fmt.Printf("logger.create.isSetted setted=%t defaultConfigIsNil=%t loggerFactoryIsNil=%t\n", setted, defaultConfig == nil, loggerFactory == nil)
		panic(ErrSetupNeverCalled)
	}
	return loggerFactory(*defaultConfig)
}

func getConfiguration(configName string) (Configuration, error) {
	var loggerConfig Configuration
	if err := config.UnmarshalKey(configName, &loggerConfig); err != nil {
		return Configuration{}, err
	}
	if loggerConfig.Debug {
		fmt.Printf("logger.getConfiguration Configuration=%s", loggerConfig.String())
	}
	return loggerConfig, nil
}

//New creates a logger implemetor with the provided configuration
func New(config Configuration) Logger {
	switch config.Provider {
	case ZAP:
		return newZap(config)
	case LOGRUS:
		return newLogrus(config)
	default:
		panic(ErrInvalidProvider)
	}
}

//NewByConfig creates a logger implemetor with the provided named configuration
func NewByConfig(configName string) Logger {
	loggerConfig, err := getConfiguration(configName)
	if err != nil {
		panic(err)
	}
	return New(loggerConfig)
}

func getOutput(out Out) (io.Writer, error) {
	switch out {
	case STDOUT:
		return os.Stdout, nil
	case STDERR:
		return os.Stderr, nil
	case DISCARD:
		return ioutil.Discard, nil
	default:
		file, err := os.OpenFile(out.String(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("CreateFileOutputErr[Out=%v Message='%v']", out, err)
		}
		return file, nil
	}
}

func Debug(message string, fields ...Field) {
	Get().Debug(message, fields...)
}

func Info(message string, fields ...Field) {
	Get().Info(message, fields...)

}

func Warn(message string, fields ...Field) {
	Get().Warn(message, fields...)
}

func Error(message string, fields ...Field) {
	Get().Error(message, fields...)
}

func Panic(message string, fields ...Field) {
	Get().Panic(message, fields...)
}

func Fatal(message string, fields ...Field) {
	Get().Fatal(message, fields...)
}

func Debugf(message string, fields ...interface{}) {
	Get().Debugf(message, fields...)
}

func Infof(message string, fields ...interface{}) {
	Get().Infof(message, fields...)
}

func Warnf(message string, fields ...interface{}) {
	Get().Warnf(message, fields...)
}

func Errorf(message string, fields ...interface{}) {
	Get().Errorf(message, fields...)
}

func Panicf(message string, fields ...interface{}) {
	Get().Panicf(message, fields...)
}

func Fatalf(message string, fields ...interface{}) {
	Get().Fatalf(message, fields...)
}

func String(key, val string) Field {
	return Field{key: key, val: val, valType: StringField}
}

func Bytes(key string, val []byte) Field {
	return Field{key: key, val: string(val), valType: BytesField}
}

func Int(key string, val int) Field {
	return Field{key: key, val: val, valType: IntField}
}

func Int64(key string, val int64) Field {
	return Field{key: key, val: val, valType: Int64Field}
}

func Float(key string, val float32) Field {
	return Field{key: key, val: val, valType: FloatField}
}

func Float64(key string, val float64) Field {
	return Field{key: key, val: val, valType: Float64Field}
}

func Bool(key string, val bool) Field {
	return Field{key: key, val: val, valType: BoolField}
}

func Duration(key string, val time.Duration) Field {
	return Field{key: key, val: val, valType: DurationField}
}

func Time(key string, val time.Time) Field {
	return Field{key: key, val: val, valType: TimeField}
}

func Struct(key string, val interface{}) Field {
	return Field{key: key, val: val, valType: StructField}
}

func Err(val error) Field {
	return Field{key: "error", val: val, valType: ErrorField}
}
