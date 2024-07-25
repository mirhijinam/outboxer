package logger

import (
	"os"

	"go.uber.org/zap"
)

func New() *zap.Logger {
	var logger *zap.Logger

	switch mode := os.Getenv("LOG_MODE"); mode {
	case "debug":
		logger = zap.Must(zap.NewProduction())
	case "develop":
		logger = zap.Must(zap.NewDevelopment())
	}

	return logger
}
