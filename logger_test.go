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

func BenchmarkSetupLogrusLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		config := &Configuration{Provider: LOGRUS}
		setupErr := Setup(config)
		assert.Nil(b, setupErr)
	}
}

func BenchmarkSetupZapLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		config := &Configuration{Provider: ZAP}
		setupErr := Setup(config)
		assert.Nil(b, setupErr)
	}
}

func BenchmarkSetupLoggingLogger(b *testing.B) {
	for k := 0; k < b.N; k++ {
		config := &Configuration{Provider: LOGGING}
		setupErr := Setup(config)
		assert.Nil(b, setupErr)
	}
}

func BenchmarkNewLogrusLogger(b *testing.B) {
	config := &Configuration{Provider: LOGRUS}
	setupErr := Setup(config)
	assert.Nil(b, setupErr)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l := NewLogger()
			assert.NotNil(b, l)
		}
	})
}

func BenchmarkNewZapLogger(b *testing.B) {
	config := &Configuration{Provider: ZAP}
	setupErr := Setup(config)
	assert.Nil(b, setupErr)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l := NewLogger()
			assert.NotNil(b, l)
		}
	})
}

func BenchmarkNewLoggingLogger(b *testing.B) {
	config := &Configuration{Provider: LOGGING}
	setupErr := Setup(config)
	assert.Nil(b, setupErr)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l := NewLogger()
			assert.NotNil(b, l)
		}
	})
}

func BenchmarkLogrusFormatfLogger(b *testing.B) {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrus.Debugf("BenchamarkLogrusFormat[field1=%v field2=%v field3=%v field4=%v field5=%v field6=%v field7=%v field8=%v field9=%v field0=%v]",
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
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrus.WithFields(logrus.Fields{
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
	setupErr := Setup(&Configuration{Provider: LOGRUS, Out: DISCARD})
	assert.Nil(b, setupErr)
	assert.Nil(b, setupErr)
	logger := NewLogger()
	assert.NotNil(b, logger)

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
	logger := zap.New(
		zap.NewTextEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	)
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
	logger := zbark.Barkify(zap.New(
		zap.NewJSONEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	))
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
	setupErr := Setup(&Configuration{Provider: ZAP, Out: DISCARD})
	assert.Nil(b, setupErr)
	logger := NewLogger()
	assert.NotNil(b, logger)

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
	backEndMessages := logging.NewBackendFormatter(logging.NewLogBackend(ioutil.Discard, "", 0), loggingFormatter)
	levelMessages := logging.AddModuleLevel(backEndMessages)
	levelMessages.SetLevel(logging.DEBUG, "")
	logging.SetBackend(levelMessages)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger := logging.MustGetLogger("benchmark")
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
	setupErr := Setup(&Configuration{Provider: LOGGING, Out: DISCARD})
	assert.Nil(b, setupErr)
	logger := NewLogger()
	assert.NotNil(b, logger)

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
