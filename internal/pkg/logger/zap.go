package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(mode, filepath string) *zap.Logger {
	cfg := zap.NewProductionConfig()

	var lvl zapcore.Level
	var encoding string
	switch mode {
	case "info":
		lvl = zap.InfoLevel
		encoding = "console"
	case "debug":
		lvl = zap.DebugLevel
		encoding = "json"
	}

	cfg.Level = zap.NewAtomicLevelAt(lvl)
	cfg.Encoding = encoding
	cfg.OutputPaths = []string{
		filepath,
	}
	cfg.ErrorOutputPaths = []string{
		filepath,
	}

	return zap.Must(cfg.Build())
}
