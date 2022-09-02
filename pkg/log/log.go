package log

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Init(v int) logr.Logger {
	zapV := -zapcore.Level(v)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		CallerKey:      "caller",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, os.Stderr, zapV)
	z := zap.New(core, zap.WithCaller(true))
	return zapr.NewLoggerWithOptions(z, zapr.LogInfoLevel("v"))
}
