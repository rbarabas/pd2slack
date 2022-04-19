package log

import (
	"fmt"

	"go.uber.org/zap"
)

func init() {
	logger, err := zap.NewProduction()
	defer logger.Sync()
	if err != nil {
		panic(fmt.Sprintf("Error initializing zap: %v+", err))
	}
	zap.ReplaceGlobals(logger)
}

func Panic(message string) {
	zap.S().Panic(message)
}

func Fatal(message string) {
	zap.S().Fatal(message)
}

func Error(message string) {
	zap.S().Error(message)
}

func Warn(message string) {
	zap.S().Warn(message)
}

func Info(message string) {
	zap.S().Info(message)
}

func Debug(message string) {
	zap.S().Debug(message)
}

func Panicf(format string, args ...interface{}) {
	zap.S().Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

func Debugf(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}
