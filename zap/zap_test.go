package zap

import (
	"errors"
	"fmt"
	"github.com/rjansen/l"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// "go.uber.org/zap/zaptest"
	"os"
	"testing"
	"time"
)

var (
	configTest  = &l.Configuration{Format: l.JSON, Out: l.DISCARD}
	myZapConfig = &l.Configuration{Format: l.JSON, Out: l.STDOUT}
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

func setupLoggerTest(t assert.TestingT, config *l.Configuration) {
	clean(t)
	setupErr := Setup(config)
	assert.Nil(t, setupErr)
}

func newLoggerTest(t assert.TestingT, config *l.Configuration) l.Logger {
	setupLoggerTest(t, config)
	logger, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	return logger
}

func newZapByConfigTest(t assert.TestingT, config *l.Configuration) l.Logger {
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
		config  l.Configuration
		success bool
	}{
		{l.Configuration{}, true},
		{l.Configuration{Out: l.STDOUT}, true},
		{l.Configuration{Out: l.STDERR}, true},
		{l.Configuration{Out: l.DISCARD}, true},
		{l.Configuration{Out: l.Out("/tmp/glive_test.log")}, true},
		{l.Configuration{Format: l.JSON}, true},
		{l.Configuration{Format: l.TEXT}, true},
		{l.Configuration{Format: l.TEXTColor}, true},
		{l.Configuration{Format: l.JSONColor}, true},
		{l.Configuration{Level: l.Level("invalid")}, true},
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

func TestNewLogger(t *testing.T) {
	setupLoggerTest(t, configTest)
	r, err := l.New()
	assert.Nil(t, err)
	assert.NotNil(t, r)
	logTest(r)
}

func TestNewZapLogger(t *testing.T) {
	l := newLoggerTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByConfig(t *testing.T) {
	l := newZapByConfigTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByInvalidLevelConfig(t *testing.T) {
	l, err := NewByConfig(&l.Configuration{Level: l.Level(99)})
	assert.Nil(t, err)
	assert.NotNil(t, l)
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

func TestNewLoggerSuccess(t *testing.T) {
	cases := []struct {
		config l.Configuration
	}{
		{l.Configuration{}},
	}
	for _, c := range cases {
		clean(t)
		setupErr := Setup(&c.config)
		assert.Nil(t, setupErr)
		logger, err := New()
		assert.Nil(t, err)
		assert.NotNil(t, logger)
		assert.NotNil(t, logger.(*zapLogger).logger)
		logger.Debug("DebugMessage", l.Struct("config", c.config))
		logger.Info("InfoMessage", l.Struct("config", c.config))
		logger.Warn("WarnMessage", l.Struct("config", c.config))
		logger.Error("ErrorMessage", l.Struct("config", c.config))
	}
}

func myZapTestSetup(t assert.TestingT) {
	clean(t)
	setupErr := Setup(myZapConfig)
	assert.Nil(t, setupErr)
}

func BenchmarkMySetupLogrusLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		myZapTestSetup(b)
	}
}

func myNewLogger(t assert.TestingT) l.Logger {
	l, err := New()
	assert.Nil(t, err)
	assert.NotNil(t, l)
	return l
}

func zapNew(t assert.TestingT) *zap.Logger {
	l := zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig),
			zapcore.AddSync(os.Stdout),
			zap.DebugLevel,
		),
	)
	assert.NotNil(t, l)
	return l
}

func zapSugarNew(t assert.TestingT) *zap.SugaredLogger {
	l := zapNew(t)
	return l.Sugar()
}

func BenchmarkNewZapLogger(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			zapNew(b)
		}
	})
}

func BenchmarkMyNewZapLogger(b *testing.B) {
	myZapTestSetup(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			myNewLogger(b)
		}
	})
}

func BenchmarkZapFormatfLogger(b *testing.B) {
	logger := zapSugarNew(b)

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

func BenchmarkZapFieldsLogger(b *testing.B) {
	logger := zapNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchamarkZapFileds",
				zap.String("field1", "field1"),
				zap.String("field2", "field2"),
				zap.String("field3", "field3"),
				zap.String("field4", "field4"),
				zap.String("field5", "field5"),
				zap.String("field6", "field6"),
				zap.String("field7", "field7"),
				zap.String("field8", "field8"),
				zap.String("field9", "field9"),
				zap.String("field0", "field0"),
			)
		}
	})
}

func BenchmarkMyZapFieldsLogger(b *testing.B) {
	myZapTestSetup(b)
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
