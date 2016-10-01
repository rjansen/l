package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestSetupLoggerSuccess(t *testing.T) {
	cases := []struct {
		config Configuration
		err    error
	}{
		{Configuration{Provider: LOGRUS}, nil},
		{Configuration{Provider: ZAP}, nil},
		{Configuration{Provider: OP, DefaultLevel: DEBUG, File: fmt.Sprintf("/tmp/op_provider_test-%d.log", rand.Int()), Format: "%{time:2006-01-02T15:04:05.999Z-07:00} %{id:03x} [%{level:.5s}] %{shortpkg}.%{longfunc} %{message}"}, nil},
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
