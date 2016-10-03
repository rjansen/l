package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/uber-common/bark"
	"github.com/uber-go/zap"
	"github.com/uber-go/zap/zbark"
	"io/ioutil"
	//"os"
	"testing"
)

var (
	myLogrusConfig     = &Configuration{Provider: LOGRUS, Out: DISCARD}
	myLogrusFmtfConfig = &Configuration{Provider: LOGRUS, Format: LOGRUSFmtfText, Out: DISCARD}
	myZapConfig        = &Configuration{Provider: ZAP, Out: DISCARD}
	myLoggingConfig    = &Configuration{Provider: LOGGING, Out: DISCARD}
)

func TestSetupLoggerSuccess(t *testing.T) {
	cases := []struct {
		config Configuration
		err    error
	}{
		{Configuration{Provider: LOGRUS}, nil},
		{Configuration{Provider: ZAP}, nil},
		{Configuration{Provider: LOGGING}, nil},
		{Configuration{Provider: Provider("invalid")}, ErrInvalidProvider},
	}
	for _, c := range cases {
		setupErr := Setup(&c.config)
		assert.Equal(t, setupErr, c.err)
	}
}

func TestSetupLoggerErrInvalidProvider(t *testing.T) {
	config := &Configuration{}
	setupErr := Setup(config)
	assert.Equal(t, setupErr, ErrInvalidProvider)
}

func TestNewLoggerSuccess(t *testing.T) {
	cases := []struct {
		config Configuration
	}{
		{Configuration{Provider: LOGRUS}},
		{Configuration{Provider: ZAP}},
		{Configuration{Provider: LOGGING}},
	}
	for _, c := range cases {
		setupErr := Setup(&c.config)
		assert.Nil(t, setupErr)
		logger := NewLogger()
		assert.NotNil(t, logger)
		switch c.config.Provider {
		case LOGRUS:
			assert.NotNil(t, logger.(*logrusLogger).logger)
		case ZAP:
			assert.NotNil(t, logger.(*zapLogger).logger)
		case LOGGING:
			assert.NotNil(t, logger.(*loggingLogger).logger)
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

func myLogrusFmtfTestSetup(t assert.TestingT) {
	setupErr := Setup(myLogrusFmtfConfig)
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

func loggingTestSetup() {
	backEndMessages := logging.NewBackendFormatter(logging.NewLogBackend(ioutil.Discard, "", 0), loggingFormatter)
	levelMessages := logging.AddModuleLevel(backEndMessages)
	levelMessages.SetLevel(logging.DEBUG, "")
	logging.SetBackend(levelMessages)
}

func myLoggingTestSetup(t assert.TestingT) {
	setupErr := Setup(myLoggingConfig)
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

func BenchmarkMySetupLogrusFmtfLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		myLogrusFmtfTestSetup(b)
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

func BenchmarkSetupLoggingLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		loggingTestSetup()
	}
}

func BenchmarkMySetupLoggingLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		myLoggingTestSetup(b)
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
		zap.NewTextEncoder(),
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

func loggingNew(t assert.TestingT) *logging.Logger {
	l := logging.MustGetLogger("benchmark")
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

func BenchmarkMyNewLogrusFmtfLogger(b *testing.B) {
	myLogrusFmtfTestSetup(b)

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

func BenchmarkNewLoggingLogger(b *testing.B) {
	loggingTestSetup()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			loggingNew(b)
		}
	})
}

func BenchmarkMyNewLoggingLogger(b *testing.B) {
	myLoggingTestSetup(b)

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

func BenchmarkMyLogrusFmtfLogger(b *testing.B) {
	myLogrusFmtfTestSetup(b)
	logger := myNewLogger(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchamarkMyLogrusFmtf",
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

func BenchmarkZapLogger(b *testing.B) {
	logger := zapNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchamarkZap",
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

func BenchmarkLoggingFormatfLogger(b *testing.B) {
	logrusTestSetup()
	logger := loggingNew(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debugf("BenchamarkLogging[field1=%v field2=%v field3=%v field4=%v field5=%v field6=%v field7=%v field8=%v field9=%v field0=%v]",
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

func BenchmarkMyLoggingFormatfLogger(b *testing.B) {
	myLoggingTestSetup(b)
	logger := myNewLogger(b)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchamarkMyLogging",
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
