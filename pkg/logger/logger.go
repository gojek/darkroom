// Package logger wraps some commonly used functions from zap.Logger and zap.SugaredLogger
// and maintains a single instance of the logger
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

// Debug logs the message at debug level with additional fields, if any
func Debug(message string, fields ...zap.Field) {
	getLogger().Debug(message, fields...)
}

// Debugf allows Sprintf style formatting and logs at debug level
func Debugf(template string, args ...interface{}) {
	getSugaredLogger().Debugf(template, args)
}

// Error logs the message at error level and prints stacktrace with additional fields, if any
func Error(message string, fields ...zap.Field) {
	getLogger().Error(message, fields...)
}

// Errorf allows Sprintf style formatting, logs at error level and prints stacktrace
func Errorf(template string, args ...interface{}) {
	getSugaredLogger().Errorf(template, args...)
}

// Fatal logs the message at fatal level with additional fields, if any and exits
func Fatal(message string, fields ...zap.Field) {
	getLogger().Fatal(message, fields...)
}

// Fatalf allows Sprintf style formatting, logs at fatal level and exits
func Fatalf(template string, args ...interface{}) {
	getSugaredLogger().Fatalf(template, args)
}

// Info logs the message at info level with additional fields, if any
func Info(message string, fields ...zap.Field) {
	getLogger().Info(message, fields...)
}

// Infof allows Sprintf style formatting and logs at info level
func Infof(template string, args ...interface{}) {
	getSugaredLogger().Infof(template, args)
}

// Warn logs the message at warn level with additional fields, if any
func Warn(message string, fields ...zap.Field) {
	getLogger().Warn(message, fields...)
}

// Warnf allows Sprintf style formatting and logs at warn level
func Warnf(template string, args ...interface{}) {
	getSugaredLogger().Warnf(template, args)
}

// AddHook adds func(zapcore.Entry) error) to the logger lifecycle
func AddHook(hook func(zapcore.Entry) error) {
	instance = getLogger().WithOptions(zap.Hooks(hook))
	sugarInstance = instance.Sugar()
}

// WithRequest takes in a http.Request and logs the message with request's Method, Host and Path
// and returns zap.logger
func WithRequest(r *http.Request) *zap.Logger {
	return getLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}

// SugaredWithRequest takes in a http.Request and logs the message with request's Method, Host and Path
// and returns zap.SugaredLogger to support Sprintf styled logging
func SugaredWithRequest(r *http.Request) *zap.SugaredLogger {
	return getSugaredLogger().With(
		zap.Any("method", r.Method),
		zap.Any("host", r.Host),
		zap.Any("path", r.URL.Path),
	)
}
