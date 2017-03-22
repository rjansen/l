package zap

import (
	"fmt"
	"github.com/rjansen/l"
)

//Hooks is the type to configure an create hooks for the logger implementation
type Hooks struct {
	Syslog SocketHook `json:"syslog" mapstructure:"syslog"`
	Gelf   SocketHook `json:"gelf" mapstructure:"gelf"`
	Stdout bool       `json:"stdout" mapstructure:"stdout"`
}

func (h Hooks) String() string {
	return fmt.Sprintf("Syslog=%s Gelf=%s Stdout=%t", h.Syslog.String(), h.Gelf.String(), h.Stdout)
}

//SocketHook is a hook that intent to sends data over network sockets
type SocketHook struct {
	Socket  string `json:"socket" mapstructure:"socket"`
	Address string `json:"addr" mapstructure:"addr"`
	Level   string `json:"level" mapstructure:"level"`
}

func (s SocketHook) String() string {
	return fmt.Sprintf("Socket=%s Address=%s Level=%s", s.Socket, s.Address, s.Level)
}

//Configuration holds the log beahvior parameters
type Configuration struct {
	Debug  bool     `json:"debug" mapstructure:"debug"`
	Level  l.Level  `json:"level" mapstructure:"level"`
	Format l.Format `json:"format" mapstructure:"format"`
	Out    l.Out    `json:"out" mapstructure:"out"`
	Hooks  Hooks    `json:"hooks" mapstructure:"hooks"`
}

func (l Configuration) String() string {
	return fmt.Sprintf("Backend=logrus Level=%s Format=%s Out=%s Hooks=%s", l.Level, l.Format, l.Out, l.Hooks)
}
