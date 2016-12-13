package migi

import (
	"time"
)

//Configuration is an interface to abstract the system parameters access
type Configuration interface {
	Debug() bool
	GetInterface(key string) interface{}
	GetBool(key string) bool
	GetDuration(key string) time.Duration
	GetFloat64(key string) float64
	GetInt(key string) int
	GetInt64(key string) int64
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetStringSlice(key string) []string
	GetTime(key string) time.Time
	InConfig(key string) bool
	IsSet(key string) bool
	//Reset()
	Unmarshal(rawVal interface{}) error
	UnmarshalKey(key string, rawVal interface{}) error
}
