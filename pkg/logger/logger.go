package logger

import (
	"net/http"
	"sync"

	"github.com/gojek/darkroom/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var instance *zap.Logger
var sugarInstance *zap.SugaredLogger
var once sync.Once

func getLogger() *zap.Logger {
	once.Do(func() {
		instance = newLogger()
	})
	return instance
}

func getSugaredLogger() *zap.SugaredLogger {
	if sugarInstance == nil {
		sugarInstance = getLogger().Sugar()
	}
	return sugarInstance
}

func getLogLevel() zapcore.Level {
	if config.LogLevel() == "" {
		return zap.DebugLevel
	}
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(config.LogLevel())); err != nil {
		Errorf("failed to parse log level from config with error: %s", err)
		panic(err)
	}
	return level
}

func newLogger() *zap.Logger {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = zap.NewAtomicLevelAt(getLogLevel())
	loggerConfig.Encoding = "json"
	loggerConfig.DisableCaller = true
	loggerConfig.OutputPaths = []string{"stdout"}
	if logger, err := loggerConfig.Build(); err != nil {
		Errorf("failed to create new logger with error: %s", err)
		panic(err)
	} else {
		return logger
	}
}

func Debug(message string, fields ...zap.Field) {
	getLogger().Debug(message, fields...)
}

func Debugf(template string, args ...interface{}) {
	getSugaredLogger().Debugf(template, args)
}

func Error(message string, fields ...zap.Field) {
	getLogger().Error(message, fields...)
}

func Errorf(template string, args ...interface{}) {
	getSugaredLogger().Errorf(template, args)
}

func Fatal(message string, fields ...zap.Field) {
	getLogger().Fatal(message, fields...)
}

func Fatalf(template string, args ...interface{}) {
	getSugaredLogger().Fatalf(template, args)
}

func Info(message string, fields ...zap.Field) {
	getLogger().Info(message, fields...)
}

func Infof(template string, args ...interface{}) {
	getSugaredLogger().Infof(template, args)
}

func Warn(message string, fields ...zap.Field) {
	getLogger().Warn(message, fields...)
}

func Warnf(template string, args ...interface{}) {
	getSugaredLogger().Warnf(template, args)
}

func AddHook(hook func(zapcore.Entry) error) {
	instance = getLogger().WithOptions(zap.Hooks(hook))
	sugarInstance = instance.Sugar()
}

func WithRequest(r *http.Request) *zap.Logger {
	return getLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}

func SugaredWithRequest(r *http.Request) *zap.SugaredLogger {
	return getSugaredLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}
