package l

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testLogger struct {
	name   string
	driver *mockDriver
	debug  testLogWriter
	info   testLogWriter
	error  testLogWriter
}

type testLogWriter struct {
	message string
	writer  *mockLogWriter
	values  []Value
}

func (scenario testLogger) setup(t *testing.T) {
	scenario.debug.writer.On("Write", mock.AnythingOfType("[]l.Value")).Once()
	scenario.driver.On("Log", DEBUG, scenario.debug.message).Return(scenario.debug.writer).Once()

	scenario.info.writer.On("Write", mock.AnythingOfType("[]l.Value")).Once()
	scenario.driver.On("Log", INFO, scenario.info.message).Return(scenario.info.writer).Once()

	scenario.error.writer.On("Write", mock.AnythingOfType("[]l.Value")).Once()
	scenario.driver.On("Log", ERROR, scenario.error.message).Return(scenario.error.writer).Once()

	scenario.driver.On("Close").Once()
}

func TestLogger(test *testing.T) {
	scenarios := []testLogger{
		{
			name:   "Creates a new Logger",
			driver: newMockDriver(),
			debug: testLogWriter{
				message: "debuglog",
				writer:  newMockLogWriter(),
				values: []Value{
					NewValue("stringvalue", "debug.stringvalue1"),
					NewValue("intvalue", 999),
					NewValue("floatvalue", 999.99),
					NewValue("timevalue", time.Now().UTC()),
					NewValue("durationvalue", time.Second*9),
					NewValue("anyvalue", new(struct{})),
				},
			},
			info: testLogWriter{
				message: "infolog",
				writer:  newMockLogWriter(),
				values: []Value{
					NewValue("stringvalue", "info.stringvalue1"),
					NewValue("intvalue", 888),
					NewValue("floatvalue", 888.88),
					NewValue("timevalue", time.Now().UTC()),
					NewValue("durationvalue", time.Second*8),
					NewValue("anyvalue", new(struct{})),
				},
			},
			error: testLogWriter{
				message: "errorlog",
				writer:  newMockLogWriter(),
				values: []Value{
					NewValue("stringvalue", "error.stringvalue1"),
					NewValue("intvalue", 777),
					NewValue("floatvalue", 777.77),
					NewValue("timevalue", time.Now().UTC()),
					NewValue("durationvalue", time.Second*7),
					NewValue("errorvalue", errors.New("errorlog")),
					NewValue("anyvalue", new(struct{})),
				},
			},
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				scenario.setup(t)

				log := New(scenario.driver)
				assert.NotNil(t, log, "logger instance")

				log.Debug(scenario.debug.message, scenario.debug.values...)
				scenario.driver.AssertCalled(t, "Log", DEBUG, scenario.debug.message)
				scenario.debug.writer.AssertCalled(t, "Write", scenario.debug.values)

				log.Info(scenario.info.message, scenario.info.values...)
				scenario.driver.AssertCalled(t, "Log", INFO, scenario.info.message)
				scenario.info.writer.AssertCalled(t, "Write", scenario.info.values)

				log.Error(scenario.error.message, scenario.error.values...)
				scenario.driver.AssertCalled(t, "Log", ERROR, scenario.error.message)
				scenario.error.writer.AssertCalled(t, "Write", scenario.error.values)

				log.Close()
			},
		)
	}
}

type testOut struct {
	name     string
	output   string
	expected Out
}

func TestOut(test *testing.T) {
	scenarios := []testOut{
		{
			name:     "Creates default STDOUT Out",
			output:   "",
			expected: STDOUT,
		},
		{
			name:     "Creates a stdout Out",
			output:   "stdout",
			expected: STDOUT,
		},
		{
			name:     "Creates a STDOUT Out",
			output:   "STDOUT",
			expected: STDOUT,
		},
		{
			name:     "Creates a file Out",
			output:   "pathtoafile",
			expected: Out("pathtoafile"),
		},
		{
			name:     "Creates a stderr Out",
			output:   "stderr",
			expected: STDERR,
		},
		{
			name:     "Creates a STDERR Out",
			output:   "STDERR",
			expected: STDERR,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				var (
					out Out
					err = out.Set(scenario.output)
				)
				assert.Nil(t, err, "Out.Set error")
				assert.Exactly(t, scenario.expected, out, "out instance")
			},
		)
	}
}

type testLevel struct {
	name     string
	level    string
	expected Level
}

func TestLevel(test *testing.T) {
	scenarios := []testLevel{
		{
			name:     "Creates default DEBUG Level",
			level:    "",
			expected: DEBUG,
		},
		{
			name:     "Creates an invalid Level",
			level:    "invalid",
			expected: DEBUG,
		},
		{
			name:     "Creates a debug Level",
			level:    "debug",
			expected: DEBUG,
		},
		{
			name:     "Creates a DEBUG Level",
			level:    "DEBUG",
			expected: DEBUG,
		},
		{
			name:     "Creates a info Level",
			level:    "info",
			expected: INFO,
		},
		{
			name:     "Creates a INFO Level",
			level:    "INFO",
			expected: INFO,
		},
		{
			name:     "Creates a error Level",
			level:    "error",
			expected: ERROR,
		},
		{
			name:     "Creates a ERROR Level",
			level:    "ERROR",
			expected: ERROR,
		},
	}

	for index, scenario := range scenarios {
		test.Run(
			fmt.Sprintf("[%d]-%s", index, scenario.name),
			func(t *testing.T) {
				var (
					level Level
					err   = level.Set(scenario.level)
				)
				assert.Nil(t, err, "Level.Set error")
				assert.Exactly(t, scenario.expected, level, "level instance")
			},
		)
	}
}
