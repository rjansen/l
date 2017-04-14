package main

import (
	"errors"
	"github.com/rjansen/l"
	"github.com/rjansen/l/zap"
	"os"
	"runtime/pprof"
	"time"
)

func init() {
	cfg := &l.Configuration{
		Level: l.DEBUG,
		Out:   l.DISCARD,
	}
	if err := zap.Setup(cfg); err != nil {
		panic(err)
	}
}

func main() {
	log, err := l.New(l.String("name", "l.zap"))
	if err != nil {
		panic(err)
	}
	out, _ := os.Create("zap-cpu.pprof")
	defer out.Close()
	defer pprof.StopCPUProfile()

	pprof.StartCPUProfile(out)
	for i := 0; i < 100000; i++ { // Arbitrary large number of iterations
		log.Debug("",
			l.String("string", "string logger field"),
			l.Bytes("bytes", []byte("[]byte logger field")),
			l.Int("int", 1),
			l.Int32("int32", 1),
			l.Int64("int64", 2),
			l.Float("float", 3.0),
			l.Float64("float", 3.0),
			l.Bool("bool", true),
			l.Duration("duration", time.Second),
			l.Time("time", time.Unix(0, 0)),
			l.Time("now", time.Now()),
			l.String("another string", "done!"),
			l.Err(errors.New("some error")),
		)
	}
}
