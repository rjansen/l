package zap

import (
	"testing"
)

func BenchmarkSetup(b *testing.B) {
	b.Run("l", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				lSetup(b, configTest)
			}
		})
	})
	// b.Run("Zap", func(b *testing.B) {
	// 	logger := newZap(b)
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			logger.Info("", fakeFields()...)
	// 		}
	// 	})
	// })
	// b.Run("Zap.Sugar", func(b *testing.B) {
	// 	logger := newZapLogger(zap.DebugLevel).Sugar()
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			logger.Infow(getMessage(0), fakeSugarFields()...)
	// 		}
	// 	})
	// })
}

func BenchmarkLog(b *testing.B) {
	b.Run("Zap", func(b *testing.B) {
		logger := zapNew(b)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(messageMock)
			}
		})
	})
	b.Run("Sugar", func(b *testing.B) {
		logger := sugarNew(b)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debugw(messageMock)
			}
		})
	})
	b.Run("l", func(b *testing.B) {
		logger := lSetupNew(b, configTest)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(messageMock)
			}
		})
	})
}

func BenchmarkLogWithFields(b *testing.B) {
	b.Run("Zap", func(b *testing.B) {
		logger := zapNew(b)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(messageMock, zapFakeFields()...)
			}
		})
	})
	b.Run("Sugar", func(b *testing.B) {
		logger := sugarNew(b)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debugw(messageMock, sugarFakeFields()...)
			}
		})
	})
	b.Run("l", func(b *testing.B) {
		logger := lSetupNew(b, configTest)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(messageMock, lFakeFields()...)
			}
		})
	})
}
