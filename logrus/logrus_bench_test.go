package logrus

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
	// b.Run("Logrus", func(b *testing.B) {
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	//			logrusSetup(b)
	// 		}
	// 	})
	// })
}

func BenchmarkLog(b *testing.B) {
	b.Run("Logrus", func(b *testing.B) {
		logger := logrusNew(b)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(messageMock)
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
	b.Run("Logrus", func(b *testing.B) {
		logger := logrusNew(b)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Debug(messageMock, logrusFakeFields())
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
