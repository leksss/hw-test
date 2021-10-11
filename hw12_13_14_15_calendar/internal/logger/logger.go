package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.Logger
}

type LoggerConf struct {
	Level string
	File  string
}

func New(config LoggerConf, projectRoot string) *Logger {
	var zapLevel zapcore.Level
	zapLevel.Set(config.Level)

	cfg := zap.NewProductionConfig()
	cfg.Level.SetLevel(zapLevel)

	cfg.OutputPaths = []string{
		"stdout",
		projectRoot + "/" + config.File,
	}

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("can't initialize logger logger: %v", err)
	}
	defer logger.Sync() // flushes buffer, if any

	return &Logger{
		logger: logger,
	}
}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}
