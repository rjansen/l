package l

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type testDriver struct {
	name     string
	logger   *mockZapLogger
	writer   *mockZapWriter
	level    Level
	zapLevel zapcore.Level
	message  string
	values   []Value
	fields   []zapcore.Field
}

func (scenario testDriver) setup(t *testing.T) {
	if scenario.writer != nil {
		scenario.writer.On("Write", mock.AnythingOfType("[]zapcore.Field")).Once()
		scenario.logger.On("Check", mock.AnythingOfType("zapcore.Level"), mock.AnythingOfType("string")).Return(
			scenario.writer,
		).Once()
	} else {
		scenario.logger.On("Check", mock.AnythingOfType("zapcore.Level"), mock.AnythingOfType("string")).Return(
			nil,
		).Once()
	}
	scenario.logger.On("Sync").Return(nil).Once()
}

func TestDriver(test *testing.T) {
	var (
		timeNow = time.Now().UTC()
		anyType = new(struct{})
	)
	scenarios := []testDriver{
		testDriver{
			name:    "Creates a new Driver and writes a debug log",
			logger:  newMockZapLogger(),
			writer:  newMockZapWriter(),
			message: "debuglog",
			level:   DEBUG,
			values: []Value{
				NewValue("stringvalue", "debug.stringvalue1"),
				NewValue("intvalue", 999),
				NewValue("intvalue32", int32(999)),
				NewValue("intvalue64", int64(999)),
				NewValue("floatvalue32", float32(999.99)),
				NewValue("floatvalue", 999.99),
				NewValue("timevalue", timeNow),
				NewValue("durationvalue", time.Second*9),
				NewValue("anyvalue", anyType),
			},
			zapLevel: zapcore.DebugLevel,
			fields: []zapcore.Field{
				zap.String("stringvalue", "debug.stringvalue1"),
				zap.Int("intvalue", 999),
				zap.Int32("intvalue32", int32(999)),
				zap.Int64("intvalue64", int64(999)),
				zap.Float32("floatvalue32", float32(999.99)),
				zap.Float64("floatvalue", 999.99),
				zap.Time("timevalue", timeNow),
				zap.Duration("durationvalue", time.Second*9),
				zap.Reflect("anyvalue", anyType),
			},
		},
		testDriver{
			name:    "Creates a new Driver and writes a info log",
			logger:  newMockZapLogger(),
			writer:  newMockZapWriter(),
			level:   INFO,
			message: "infolog",
			values: []Value{
				NewValue("stringvalue", "info.stringvalue1"),
				NewValue("intvalue", 999),
				NewValue("floatvalue", 999.99),
				NewValue("timevalue", timeNow),
				NewValue("durationvalue", time.Second*9),
				NewValue("anyvalue", anyType),
			},
			zapLevel: zapcore.InfoLevel,
			fields: []zapcore.Field{
				zap.String("stringvalue", "info.stringvalue1"),
				zap.Int64("intvalue", 999),
				zap.Float64("floatvalue", 999.99),
				zap.Time("timevalue", timeNow),
				zap.Duration("durationvalue", time.Second*9),
				zap.Reflect("anyvalue", anyType),
			},
		},
		testDriver{
			name:    "Creates a new Driver and writes a error log",
			logger:  newMockZapLogger(),
			writer:  newMockZapWriter(),
			level:   ERROR,
			message: "errorlog",
			values: []Value{
				NewValue("stringvalue", "error.stringvalue1"),
				NewValue("intvalue", 999),
				NewValue("floatvalue", 999.99),
				NewValue("timevalue", timeNow),
				NewValue("durationvalue", time.Second*9),
				NewValue("errorvalue", errors.New("errorlog")),
				NewValue("anyvalue", anyType),
			},
			zapLevel: zapcore.ErrorLevel,
			fields: []zapcore.Field{
				zap.String("stringvalue", "error.stringvalue1"),
				zap.Int64("intvalue", 999),
				zap.Float64("floatvalue", 999.99),
				zap.Time("timevalue", timeNow),
				zap.Duration("durationvalue", time.Second*9),
				zap.NamedError("errorvalue", errors.New("errorlog")),
				zap.Reflect("anyvalue", anyType),
			},
		},
		testDriver{
			name:    "Creates a new Driver but does not provide any writer to log message",
			logger:  newMockZapLogger(),
			writer:  nil,
			level:   DEBUG,
			message: "debugdisabledlog",
			values: []Value{
				NewValue("stringvalue", "debugdisabled.stringvalue1"),
			},
			zapLevel: zapcore.DebugLevel,
			fields: []zapcore.Field{
				zap.String("stringvalue", "debugdisabled.stringvalue1"),
			},
		},
		testDriver{
			name:    "Creates a new Driver but does not provide any writer to an invalid level",
			logger:  newMockZapLogger(),
			writer:  nil,
			level:   Level("invalid"),
			message: "invalidlog",
			values: []Value{
				NewValue("stringvalue", "debugdisabled.stringvalue1"),
			},
			zapLevel: zapcore.DebugLevel,
			fields: []zapcore.Field{
				zap.String("stringvalue", "debugdisabled.stringvalue1"),
			},
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				var driver Driver = NewDriver(scenario.logger)
				assert.NotNil(t, driver, "driver instance")

				writer := driver.Log(scenario.level, scenario.message)
				if scenario.level == DEBUG || scenario.level == INFO || scenario.level == ERROR {
					scenario.logger.AssertCalled(t, "Check", scenario.zapLevel, scenario.message)
				} else {
					scenario.logger.AssertNotCalled(t, "Check", scenario.zapLevel, scenario.message)
				}
				if scenario.writer != nil {
					assert.NotNil(t, writer, "writer instance")
					writer.Write(scenario.values...)
					scenario.writer.AssertCalled(t, "Write", scenario.fields)
				}
				driver.Close()
				scenario.logger.AssertCalled(t, "Sync")
			},
		)
	}
}

type testZapLoggerDelegate struct {
	name    string
	zapMock *testZapMock
	level   zapcore.Level
	message string
	fields  []zapcore.Field
}

type testZapMock struct {
	logger   *zap.Logger
	core     zapcore.Core
	observer *observer.ObservedLogs
}

func newTestZapMock(level zapcore.Level) *testZapMock {
	core, observer := observer.New(level)
	return &testZapMock{
		logger:   zap.New(core),
		core:     core,
		observer: observer,
	}
}

func (scenario testZapLoggerDelegate) setup(t *testing.T) {
}

func TestZapLoggerDelegate(test *testing.T) {
	scenarios := []testZapLoggerDelegate{
		testZapLoggerDelegate{
			name:    "Creates a new zapLogger and writes a debug log",
			zapMock: newTestZapMock(zapcore.DebugLevel),
			level:   zapcore.DebugLevel,
			message: "debuglog",
			fields: []zapcore.Field{
				zap.String("stringfield", "debug.stringvalue1"),
				zap.Int64("intfield", 777),
				zap.Float64("floatfield", 777.77),
				zap.Reflect("reflectfield", new(struct{})),
				zap.Time("timefield", time.Now().UTC()),
				zap.Duration("durationfield", time.Second*7),
				zap.NamedError("errorfield", errors.New("errorlog")),
			},
		},
		testZapLoggerDelegate{
			name:    "Creates a new zapLogger and writes a info log",
			zapMock: newTestZapMock(zapcore.InfoLevel),
			level:   zapcore.InfoLevel,
			message: "infolog",
			fields: []zapcore.Field{
				zap.String("stringfield", "info.stringvalue1"),
				zap.Int64("intfield", 888),
				zap.Float64("floatfield", 888.88),
				zap.Reflect("reflectfield", new(struct{})),
				zap.Time("timefield", time.Now().UTC()),
				zap.Duration("durationfield", time.Second*8),
				zap.NamedError("errorfield", errors.New("errorlog")),
			},
		},
		testZapLoggerDelegate{
			name:    "Creates a new zapLogger and writes a error log",
			zapMock: newTestZapMock(zapcore.ErrorLevel),
			level:   zapcore.ErrorLevel,
			message: "errorlog",
			fields: []zapcore.Field{
				zap.String("stringfield", "error.stringvalue1"),
				zap.Int64("intfield", 666),
				zap.Float64("floatfield", 666.66),
				zap.Reflect("reflectfield", new(struct{})),
				zap.Time("timefield", time.Now().UTC()),
				zap.Duration("durationfield", time.Second*6),
				zap.NamedError("errorfield", errors.New("errorlog")),
			},
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				logger := newZapLoggerDelegate(scenario.zapMock.logger)
				assert.NotNil(t, logger, "delegateLogger instance")

				writer := logger.Check(scenario.level, scenario.message)
				assert.NotNil(t, writer, "writer instance")
				writer.Write(scenario.fields...)
				observedLogs := scenario.zapMock.observer.All()
				assert.NotEmpty(t, observedLogs, "observed logs")
				for _, logEntry := range observedLogs {
					assert.Equal(t, scenario.level, logEntry.Level, "log level")
					assert.Equal(t, scenario.message, logEntry.Message, "log message")
					assert.Equal(t, scenario.fields, logEntry.Context)
				}
				assert.NoError(t, logger.Sync(), "deleagteLogger sync")
			},
		)
	}
}

type testZapLogger struct {
	name   string
	output Out
	level  Level
	err    error
}

func (scenario testZapLogger) setup(t *testing.T) {
}

func TestZapLogger(test *testing.T) {
	scenarios := []testZapLogger{
		testZapLogger{
			name:   "Creates a new debug level zap logger instance",
			output: STDOUT,
			level:  DEBUG,
		},
		testZapLogger{
			name:   "Creates a new info level zap logger instance",
			output: STDOUT,
			level:  INFO,
		},
		testZapLogger{
			name:   "Creates a new error level zap logger instance",
			output: STDOUT,
			level:  ERROR,
		},
		testZapLogger{
			name:   "Creates a new zap logger instance with stderr output",
			output: STDERR,
			level:  DEBUG,
		},
		testZapLogger{
			name:   "Does not creates a new zap logger with invalid level",
			output: STDOUT,
			level:  Level("invalid"),
			err:    errors.New("unrecognized level: \"invalid\""),
		},
		testZapLogger{
			name:   "Does not creates a new zap logger with invalid output",
			output: Out(""),
			level:  DEBUG,
			err:    errors.New("couldn't open sink \"\": open : no such file or directory"),
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				var logger, err = NewZapLogger(scenario.level, scenario.output)
				assert.Equal(t, scenario.err, err, "error instance")
				if scenario.err == nil {
					assert.NotNil(t, logger, "zap.Logger instance")
				} else {
					assert.Nil(t, logger, "zap.Logger instance")
				}
			},
		)
	}
}

func TestZapLoggerDefault(test *testing.T) {
	logger := NewZapLoggerDefault()
	assert.NotNil(test, logger, "loggerDefault instance")
}
