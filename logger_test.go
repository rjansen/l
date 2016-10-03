package logger

import (
	"github.com/Sirupsen/logrus"
	"github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
	"github.com/uber-go/zap"
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
		logger.Debug("DebugMessage", Field{key: "config", val: c.config})
		logger.Info("InfoMessage", Field{key: "config", val: c.config})
		logger.Warn("WarnMessage", Field{key: "config", val: c.config})
		logger.Error("ErrorMessage", Field{key: "config", val: c.config})
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

func BenchmarkLogrusFromatfLogger(b *testing.B) {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrus.Debugf("BenchaMarkLogrusFormat[field1=%v field2=%v field3=%v field4=%v field5=%v field6=%v field7=%v field8=%v field9=%v field0=%v]",
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
			}).Debug("BenchaMarkLogrusFileds")
		}
	})
}

func BenchmarkZapLogger(b *testing.B) {
	logger := zap.New(
		zap.NewJSONEncoder(),
		zap.DebugLevel,
		zap.DiscardOutput,
	)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debug("BenchaMarkZap",
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

func BenchmarkLoggingFromatfLogger(b *testing.B) {
	backEndMessages := logging.NewBackendFormatter(logging.NewLogBackend(ioutil.Discard, "", 0), loggingFormatter)
	levelMessages := logging.AddModuleLevel(backEndMessages)
	levelMessages.SetLevel(logging.DEBUG, "")
	logging.SetBackend(levelMessages)

	logger := logging.MustGetLogger("benchmark")
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Debugf("BenchaMarkLogrusFormat[field1=%v field2=%v field3=%v field4=%v field5=%v field6=%v field7=%v field8=%v field9=%v field0=%v]",
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
