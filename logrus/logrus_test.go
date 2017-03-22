package logrus

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/rjansen/l"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
	"time"
)

var (
	configTest     = &Configuration{Format: l.TEXT, Out: l.DISCARD}
	myLogrusConfig = &Configuration{Format: l.TEXT, Out: l.DISCARD}
)

func clean(t assert.TestingT) {
	defaultConfig = nil
	// loggerFactory = nil
	reset(t)
}

func reset(t assert.TestingT) {
	// rootLogger = nil
	// once.Reset()
}

func setupLoggerTest(t assert.TestingT, config *Configuration) {
	clean(t)
	setupErr := Setup(config)
	assert.Nil(t, setupErr)
}

func newLoggerTest(t assert.TestingT, config *Configuration) l.Logger {
	setupLoggerTest(t, config)
	logger := New()
	assert.NotNil(t, logger)
	return logger
}

func newLogrusByConfigTest(t assert.TestingT, config *Configuration) l.Logger {
	logger, err := NewByConfig(config)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	return logger
}

func logTest(logger l.Logger) {
	logger.Debug("DebugMessage")
	logger.Debug("DebugFieldsMessage",
		l.String("Param1", "1"),
		l.String("Param2", "2"),
		l.String("Param3", "3"),
	)
	logger.Info("InfoMessage")
	logger.Info("InfoFieldsMessage",
		l.String("Param1", "1"),
		l.String("Param2", "2"),
		l.String("Param3", "3"),
	)
	logger.Warn("WarnMessage")
	logger.Warn("WarnFieldsMessage",
		l.String("Param1", "1"),
		l.String("Param2", "2"),
		l.String("Param3", "3"),
	)
	logger.Error("ErrorMessage")
	logger.Error("ErrorFieldsMessage",
		l.String("Param1", "1"),
		l.String("Param2", "2"),
		l.String("Param3", "3"),
	)
}

func TestSetupLogger(t *testing.T) {
	cases := []struct {
		config  Configuration
		success bool
	}{
		{Configuration{}, true},
		{Configuration{Out: l.STDOUT}, true},
		{Configuration{Out: l.STDERR}, true},
		{Configuration{Out: l.DISCARD}, true},
		{Configuration{Out: l.Out("/tmp/glive_test.log")}, true},
		{Configuration{Format: l.JSON}, true},
		{Configuration{Format: l.TEXT}, true},
		{Configuration{Format: l.TEXTColor}, true},
		{Configuration{Format: l.JSONColor}, true},
		{Configuration{Level: l.Level("invalid")}, false},
	}
	for _, c := range cases {
		clean(t)
		err := Setup(&c.config)
		if c.success {
			assert.Nil(t, err, fmt.Sprintf("Not Nil for data=%+v", c))
		} else {
			assert.NotNil(t, err, fmt.Sprintf("Nil for data=%+v", c))
		}
	}
}

func TestGetLogger(t *testing.T) {
	setupLoggerTest(t, configTest)
	r := l.Get()
	assert.NotNil(t, r)
	logTest(r)
}

func TestNewLogger(t *testing.T) {
	setupLoggerTest(t, configTest)
	r := l.New()
	assert.NotNil(t, r)
	logTest(r)
}

func TestNewLogrusLogger(t *testing.T) {
	l := newLoggerTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByConfig(t *testing.T) {
	l := newLogrusByConfigTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByInvalidLevelConfig(t *testing.T) {
	l, err := NewByConfig(&Configuration{Level: l.Level(99)})
	assert.NotNil(t, err)
	assert.Nil(t, l)
}

func TestOutString(t *testing.T) {
	cases := []struct {
		out l.Out
	}{
		{l.STDOUT},
		{l.STDERR},
		{l.DISCARD},
		//{OUT("invalid"), false},
	}
	for _, c := range cases {
		assert.NotNil(t, c.out.String())
	}
}

func TestOutSet(t *testing.T) {
	originalValue := "originalValue"
	var o l.Out
	assert.Nil(t, o.Set(originalValue))
	assert.Equal(t, originalValue, o.String())
}

func TestFormatString(t *testing.T) {
	cases := []struct {
		format l.Format
	}{
		{l.JSON},
		{l.TEXT},
		{l.Format("invalid")},
	}
	for _, c := range cases {
		assert.NotNil(t, c.format.String())
	}
}

func TestFormatSet(t *testing.T) {
	originalValue := "originalValue"
	var f l.Format
	assert.Nil(t, f.Set(originalValue))
	assert.Equal(t, originalValue, f.String())
}

func TestSetupLoggerErrInvalidProvider(t *testing.T) {
	clean(t)
	config := &Configuration{}
	assert.Panics(t, func() {
		NewByConfig(config)
	})
}

func TestNewLoggerSuccess(t *testing.T) {
	cases := []struct {
		config Configuration
	}{
		{Configuration{}},
	}
	for _, c := range cases {
		clean(t)
		setupErr := Setup(&c.config)
		assert.Nil(t, setupErr)
		logger := New()
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.(*logrusLogger).logger)
		logger.Debug("DebugMessage", l.Struct("config", c.config))
		logger.Info("InfoMessage", l.Struct("config", c.config))
		logger.Warn("WarnMessage", l.Struct("config", c.config))
		logger.Error("ErrorMessage", l.Struct("config", c.config))
	}
}

func logrusTestSetup() {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)
}

func myLogrusTestSetup(t assert.TestingT) {
	clean(t)
	setupErr := Setup(myLogrusConfig)
	assert.Nil(t, setupErr)
}

func BenchmarkSetupLogrusLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		logrusTestSetup()
	}
}

func BenchmarkMySetupLogrusLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		myLogrusTestSetup(b)
	}
}

func myNewLogger(t assert.TestingT) l.Logger {
	l := New()
	assert.NotNil(t, l)
	return l
}

func logrusNew(t assert.TestingT) *logrus.Logger {
	l := logrus.New()
	assert.NotNil(t, l)
	return l
}

func BenchmarkNewLogrusLogger(b *testing.B) {
	logrusTestSetup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrusNew(b)
		}
	})
}

func BenchmarkMyNewLogrusLogger(b *testing.B) {
	myLogrusTestSetup(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			myNewLogger(b)
		}
	})
}

func BenchmarkLogrusFormatfLogger(b *testing.B) {
	logrusTestSetup()
	logger := logrusNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debugf("BenchamarkLogrusFormat[field1=%v field2=%v field3=%v field4=%v field5=%v field6=%v field7=%v field8=%v field9=%v field0=%v]",
				"field1",
				"field2",
				"field3",
				"field4",
				"field5",
				"field6",
				"field7",
				"field8",
				"field9",
				"field0",
			)
		}
	})
}

func BenchmarkLogrusFieldsLogger(b *testing.B) {
	logrusTestSetup()
	logger := logrusNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(logrus.Fields{
				"field1": "field1",
				"field2": "field2",
				"field3": "field3",
				"field4": "field4",
				"field5": "field5",
				"field6": "field6",
				"field7": "field7",
				"field8": "field8",
				"field9": "field9",
				"field0": "field0",
			}).Debug("BenchamarkLogrusFileds")
		}
	})
}

func BenchmarkMyLogrusFieldsLogger(b *testing.B) {
	myLogrusTestSetup(b)
	logger := myNewLogger(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchamarkMyLogrusFileds",
				l.String("field1", "field1"),
				l.String("field2", "field2"),
				l.String("field3", "field3"),
				l.String("field4", "field4"),
				l.String("field5", "field5"),
				l.String("field6", "field6"),
				l.String("field7", "field7"),
				l.String("field8", "field8"),
				l.String("field9", "field9"),
				l.String("field0", "field0"),
			)
		}
	})
}

var errExample = errors.New("fail")

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var _jane = user{
	Name:      "Jane Doe",
	Email:     "jane@test.com",
	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
}
