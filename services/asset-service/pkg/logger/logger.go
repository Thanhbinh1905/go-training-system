package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger tạo một logger mới với file log chỉ định và service name để tag vào log.
func NewLogger(logFilePath string, serviceName string) *zap.Logger {
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// Tạo folder nếu chưa có
	if err := os.MkdirAll("logs", os.ModePerm); err != nil {
		log.Fatalf("Failed to create log folder: %v", err)
	}

	// Mở file log
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Ghi log ra cả console và file
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), zapcore.InfoLevel),
	)

	// Tạo logger kèm theo field "service"
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).With(zap.String("service", serviceName))
	return logger
}
