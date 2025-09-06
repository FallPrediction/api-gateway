package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger
var loggerOnce sync.Once

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/api-gateway.log",
		MaxSize:    500,
		MaxBackups: 2,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func NewLogger() *zap.Logger {
	if logger == nil {
		loggerOnce.Do(func() {
			logLevel, err := zapcore.ParseLevel(os.Getenv("LOG_LEVEL"))
			if err != nil {
				logLevel = zapcore.ErrorLevel
			}
			core := zapcore.NewCore(getEncoder(), getLogWriter(), logLevel)
			logger = zap.New(core, zap.AddStacktrace(logLevel))
		})
	}
	return logger
}
