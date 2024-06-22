package logger

import (
	"go-template/pkg/common/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func Initialize() error {
	encoding := "console"
	encodeLevel := zapcore.CapitalColorLevelEncoder
	if env.MODE == "production" {
		encoding = "json"
		encodeLevel = zapcore.CapitalLevelEncoder
	}

	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         encoding,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

func GetInstance() *zap.Logger {
	if logger == nil {
		once.Do(func() {
			if err := Initialize(); err != nil {
				panic(err)
			}
		})
	}

	return logger
}

func NoOpLogger() *zap.Logger {
	return zap.NewNop()
}
