package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/uber-common/bark"
	"github.com/uber-go/zap"
	"github.com/uber-go/zap/zbark"
	"io/ioutil"
	//"os"
	"errors"
	"fmt"
	"testing"
	"time"
)

var (
	configTest     = Configuration{Provider: LOGRUS, Format: TEXT, Out: DISCARD}
	myLogrusConfig = Configuration{Provider: LOGRUS, Format: TEXT, Out: DISCARD}
	myZapConfig    = Configuration{Provider: ZAP, Out: DISCARD}
)

func setupLoggerTest(t assert.TestingT, config Configuration) {
	setupErr := Setup(config)
	assert.Nil(t, setupErr)
}

func newLoggerTest(t assert.TestingT, config Configuration) Logger {
	setupLoggerTest(t, config)
	logger := NewLogger()
	assert.NotNil(t, logger)
	return logger
}

func newLogrusByConfigTest(t assert.TestingT, config Configuration) Logger {
	logger := NewLoggerByConfig(config)
	assert.NotNil(t, logger)
	return logger
}

func logTest(logger Logger) {
	logger.Debug("DebugMessage")
	logger.Debug("DebugFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
	logger.Info("InfoMessage")
	logger.Info("InfoFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
	logger.Warn("WarnMessage")
	logger.Warn("WarnFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
	logger.Error("ErrorMessage")
	logger.Error("ErrorFieldsMessage",
		String("Param1", "1"),
		String("Param2", "2"),
		String("Param3", "3"),
	)
}

func TestSetupLogger(t *testing.T) {
	cases := []struct {
		config  Configuration
		success bool
	}{
		{Configuration{Out: STDOUT}, false},
		{Configuration{Out: STDERR}, false},
		{Configuration{Out: DISCARD}, false},
		{Configuration{Provider: LOGRUS}, true},
		{Configuration{Provider: ZAP}, true},
		{Configuration{Provider: LOGRUS, Out: STDOUT}, true},
		{Configuration{Provider: ZAP, Out: STDOUT}, true},
		{Configuration{Provider: LOGRUS, Out: STDERR}, true},
		{Configuration{Provider: ZAP, Out: STDERR}, true},
		{Configuration{Provider: LOGRUS, Out: DISCARD}, true},
		{Configuration{Provider: ZAP, Out: DISCARD}, true},
		{Configuration{Provider: LOGRUS, Out: Out("/tmp/glive_test.log")}, true},
		{Configuration{Provider: ZAP, Out: Out("/tmp/glive_test.log")}, true},
		//{Configuration{Provider: LOGRUS, Out: Out("%$$/{}")}, false},
		//{Configuration{Provider: ZAP, Out: Out("$$%/{}")}, false},
		{Configuration{Provider: LOGRUS, Format: JSON}, true},
		{Configuration{Provider: ZAP, Format: JSON}, true},
		{Configuration{Provider: LOGRUS, Format: TEXT}, true},
		{Configuration{Provider: ZAP, Format: TEXT}, true},
		{Configuration{Provider: LOGRUS, Format: TEXTColor}, true},
		{Configuration{Provider: ZAP, Format: TEXTColor}, true},
		{Configuration{Provider: LOGRUS, Format: JSONColor}, true},
		{Configuration{Provider: ZAP, Format: JSONColor}, true},
		{Configuration{Provider: LOGRUS, Level: Level("invalid")}, false},
		{Configuration{Provider: ZAP, Level: Level("notvalid")}, true},
		//{Configuration{Provider: LOGRUS, Level: Level(99)}, false},
		//{Configuration{Provider: ZAP, Level: Level(99)}, false},
	}
	for _, c := range cases {
		err := Setup(c.config)
		if c.success {
			assert.Nil(t, err, fmt.Sprintf("Not Nil for data=%+v", c))
		} else {
			assert.NotNil(t, err, fmt.Sprintf("Nil for data=%+v", c))
		}
	}
}

func TestGetLogger(t *testing.T) {
	setupLoggerTest(t, configTest)
	r := GetLogger()
	assert.NotNil(t, r)
	logTest(r)
}

func TestNewLogger(t *testing.T) {
	l := newLoggerTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByConfig(t *testing.T) {
	l := newLogrusByConfigTest(t, configTest)
	logTest(l)
}

func TestNewLoggerByInvalidLevelConfig(t *testing.T) {
	l := NewLoggerByConfig(Configuration{Level: Level(99)})
	assert.NotNil(t, l)
}

func TestOutString(t *testing.T) {
	cases := []struct {
		out Out
	}{
		{STDOUT},
		{STDERR},
		{DISCARD},
		//{OUT("invalid"), false},
	}
	for _, c := range cases {
		assert.NotNil(t, c.out.String())
	}
}

func TestOutSet(t *testing.T) {
	originalValue := "originalValue"
	var o Out
	assert.Nil(t, o.Set(originalValue))
	assert.Equal(t, originalValue, o.String())
}

func TestFormatString(t *testing.T) {
	cases := []struct {
		format Format
	}{
		{JSON},
		{TEXT},
		{Format("invalid")},
	}
	for _, c := range cases {
		assert.NotNil(t, c.format.String())
	}
}

func TestFormatSet(t *testing.T) {
	originalValue := "originalValue"
	var f Format
	assert.Nil(t, f.Set(originalValue))
	assert.Equal(t, originalValue, f.String())
}

func TestSetupLoggerErrInvalidProvider(t *testing.T) {
	config := Configuration{}
	setupErr := Setup(config)
	assert.Equal(t, setupErr, ErrInvalidProvider)
}

func TestNewLoggerSuccess(t *testing.T) {
	cases := []struct {
		config Configuration
	}{
		{Configuration{Provider: LOGRUS}},
		{Configuration{Provider: ZAP}},
	}
	for _, c := range cases {
		setupErr := Setup(c.config)
		assert.Nil(t, setupErr)
		logger := NewLogger()
		assert.NotNil(t, logger)
		switch c.config.Provider {
		case LOGRUS:
			assert.NotNil(t, logger.(*logrusLogger).logger)
		default:
			assert.NotNil(t, logger.(*zapLogger).logger)
		}
		logger.Debug("DebugMessage", Struct("config", c.config))
		logger.Info("InfoMessage", Struct("config", c.config))
		logger.Warn("WarnMessage", Struct("config", c.config))
		logger.Error("ErrorMessage", Struct("config", c.config))
	}
}

func logrusTestSetup() {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)
}

func myLogrusTestSetup(t assert.TestingT) {
	setupErr := Setup(myLogrusConfig)
	assert.Nil(t, setupErr)
}

func zapTestSetup(t assert.TestingT) {
	logger := zap.New(
		zap.NewTextEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	)
	assert.NotNil(t, logger)
}

func barkifyZapTestSetup(t assert.TestingT) {
	logger := zbark.Barkify(zap.New(
		zap.NewJSONEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	))
	assert.NotNil(t, logger)
}

func myZapTestSetup(t assert.TestingT) {
	setupErr := Setup(myZapConfig)
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

func BenchmarkSetupZapLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		zapTestSetup(b)
	}
}
func BenchmarkSetupBarkifyZapLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		barkifyZapTestSetup(b)
	}
}

func BenchmarkMySetupZapLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		myZapTestSetup(b)
	}
}

func myNewLogger(t assert.TestingT) Logger {
	l := NewLogger()
	assert.NotNil(t, l)
	return l
}

func logrusNew(t assert.TestingT) *logrus.Logger {
	l := logrus.New()
	assert.NotNil(t, l)
	return l
}

func zapNew(t assert.TestingT) zap.Logger {
	l := zap.New(
		zap.NewJSONEncoder(),
		//zap.NewTextEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	)
	assert.NotNil(t, l)
	return l
}

func barkifyZapNew(t assert.TestingT) bark.Logger {
	l := zbark.Barkify(zap.New(
		zap.NewJSONEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	))
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

func BenchmarkNewZapLogger(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			zapNew(b)
		}
	})
}

func BenchmarkNewBarkifyZapLogger(b *testing.B) {
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			barkifyZapNew(b)
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
				String("field1", "field1"),
				String("field2", "field2"),
				String("field3", "field3"),
				String("field4", "field4"),
				String("field5", "field5"),
				String("field6", "field6"),
				String("field7", "field7"),
				String("field8", "field8"),
				String("field9", "field9"),
				String("field0", "field0"),
			)
		}
	})
}

func createZapFields() []zap.Field {
	return []zap.Field{
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
	}
}

func BenchmarkZapLogger(b *testing.B) {
	logger := zapNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("BenchamarkZap", createZapFields()...)
		}
	})
}

var errExample = errors.New("fail")

type user struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func (u user) MarshalLog(kv zap.KeyValue) error {
	kv.AddString("name", u.Name)
	kv.AddString("email", u.Email)
	kv.AddInt64("created_at", u.CreatedAt.UnixNano())
	return nil
}

var _jane = user{
	Name:      "Jane Doe",
	Email:     "jane@test.com",
	CreatedAt: time.Date(1980, 1, 1, 12, 0, 0, 0, time.UTC),
}

func fakeFields() []zap.Field {
	return []zap.Field{
		zap.Int("int", 1),
		zap.Int64("int64", 2),
		zap.Float64("float", 3.0),
		zap.String("string", "four!"),
		zap.Bool("bool", true),
		zap.Time("time", time.Unix(0, 0)),
		zap.Error(errExample),
		zap.Duration("duration", time.Second),
		zap.Marshaler("user-defined type", _jane),
		zap.String("another string", "done!"),
	}
}

func BenchmarkZapAddingFields(b *testing.B) {
	logger := zap.New(
		zap.NewJSONEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Go fast.", fakeFields()...)
		}
	})
}

func BenchmarkZapBarkifyLogger(b *testing.B) {
	logger := barkifyZapNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.WithFields(bark.Fields{
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
			}).Info("BenchamarkBarkifiedZap")
		}
	})
}

func BenchmarkMyZapLogger(b *testing.B) {
	myZapTestSetup(b)
	logger := myNewLogger(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchamarkMyZap",
				String("field1", "field1"),
				String("field2", "field2"),
				String("field3", "field3"),
				String("field4", "field4"),
				String("field5", "field5"),
				String("field6", "field6"),
				String("field7", "field7"),
				String("field8", "field8"),
				String("field9", "field9"),
				String("field0", "field0"),
			)
		}
	})
}
