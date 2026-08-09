package logger

import (
	"os"
	"slices"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var loggerOnce sync.Once

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func NewLogger() *zap.Logger {
	if logger == nil {
		loggerOnce.Do(func() {
			stdoutLevelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return l == zap.DebugLevel || l == zap.InfoLevel
			})
			stderrLevelEnabler := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
				return slices.Contains([]zapcore.Level{
					zap.ErrorLevel,
					zap.DPanicLevel,
					zap.PanicLevel,
				}, l)
			})
			core := zapcore.NewTee(
				zapcore.NewCore(getEncoder(), zapcore.Lock(os.Stdout), stdoutLevelEnabler),
				zapcore.NewCore(getEncoder(), zapcore.Lock(os.Stderr), stderrLevelEnabler),
			)
			logger = zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel))
		})
	}
	return logger
}
