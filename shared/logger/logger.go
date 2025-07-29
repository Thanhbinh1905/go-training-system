package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(isDev bool) {
	var err error

	if isDev {
		Log, err = zap.NewDevelopment()
		if err != nil {
			panic("failed to init dev logger: " + err.Error())
		}
		return
	}

	logPath := "./logs/app.log"
	_ = os.MkdirAll("./logs", os.ModePerm)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("failed to open log file: " + err.Error())
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(file),
		zapcore.InfoLevel,
	)

	Log = zap.New(core, zap.AddCaller())
}
