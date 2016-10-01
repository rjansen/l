package logger

import (
	logrus "github.com/Sirupsen/logrus"
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
		{Configuration{Provider: OP}, nil},
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

func BenchmarkSetupLogrusLoggerSuccess(b *testing.B) {
	for k := 0; k < b.N; k++ {
		config := &Configuration{Provider: LOGRUS}
		setupErr := Setup(config)
		assert.Nil(b, setupErr)
	}
}

func BenchmarkSetupZapLoggerSuccess(b *testing.B) {
	for k := 0; k < b.N; k++ {
		config := &Configuration{Provider: ZAP}
		setupErr := Setup(config)
		assert.Nil(b, setupErr)
	}
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
