package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(serviceName, logFilePath string) *zap.Logger {
	// Encoder config (JSON format)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Multi-output: console + file
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(&lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    10, // megabytes
			MaxBackups: 5,
			MaxAge:     30, // days
			Compress:   true,
		}), zapcore.InfoLevel),
	)

	// Create logger with service field
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).With(
		zap.String("service", serviceName),
		zap.String("env", os.Getenv("ENV")),
		zap.Time("start_time", time.Now()),
	)

	return logger
}
