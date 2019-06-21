package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"testing"
)

func makeLogsStorageHooks() (func(zapcore.Entry) error, **[]zapcore.Entry) {
	var pointerToStorage *[]zapcore.Entry
	initialStorage := make([]zapcore.Entry, 0)
	pointerToStorage = &initialStorage
	return func(entry zapcore.Entry) error {
		newStorage := append(*pointerToStorage, entry)
		pointerToStorage = &newStorage
		return nil
	}, &pointerToStorage
}

func TestLogger(t *testing.T) {
	hook, ptr := makeLogsStorageHooks()
	AddHook(hook)

	cases := []struct {
		message  string
		level    zapcore.Level
		callFunc func(message string, fields ...zap.Field)
	}{
		{
			message:  "warning message",
			level:    zap.WarnLevel,
			callFunc: Warn,
		},
		{
			message:  "debug message",
			level:    zap.DebugLevel,
			callFunc: Debug,
		},
		{
			message:  "error message",
			level:    zap.ErrorLevel,
			callFunc: Error,
		},
		{
			message:  "info message",
			level:    zap.InfoLevel,
			callFunc: Info,
		},
	}

	for i, c := range cases {
		c.callFunc(c.message)
		logsStorage := **ptr
		assert.Equal(t, i+1, len(logsStorage))
		assert.Equal(t, c.level, logsStorage[len(logsStorage)-1].Level)
		assert.Equal(t, c.message, logsStorage[len(logsStorage)-1].Message)
	}
}

func TestLoggerFormatter(t *testing.T) {
	hook, ptr := makeLogsStorageHooks()
	AddHook(hook)

	cases := []struct {
		template string
		level    zapcore.Level
		args     []string
		callFunc func(format string, args ...interface{})
	}{
		{
			template: "warning message: %s",
			level:    zap.WarnLevel,
			args:     []string{"WARNING"},
			callFunc: Warnf,
		},
		{
			template: "debug message: %s",
			level:    zap.DebugLevel,
			args:     []string{"DEBUG"},
			callFunc: Debugf,
		},
		{
			template: "error message: %s",
			level:    zap.ErrorLevel,
			args:     []string{"ERROR"},
			callFunc: Errorf,
		},
		{
			template: "info message: %s",
			level:    zap.InfoLevel,
			args:     []string{"INFO"},
			callFunc: Infof,
		},
	}

	for i, c := range cases {
		c.callFunc(c.template, c.args)
		logsStorage := **ptr
		assert.Equal(t, i+1, len(logsStorage))
		assert.Equal(t, c.level, logsStorage[len(logsStorage)-1].Level)
	}
}

func TestLoggerWithRequest(t *testing.T) {
	hook, pointerToLogsStorage := makeLogsStorageHooks()
	AddHook(hook)
	msg := "success"
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/path/to/resource", nil)

	WithRequest(req).Info(msg)

	storage := **pointerToLogsStorage
	assert.Equal(t, 1, len(storage))
	assert.Equal(t, zap.InfoLevel, storage[0].Level)
	assert.Equal(t, msg, storage[0].Message)
}

func TestSugaredLoggerWithRequest(t *testing.T) {
	hook, pointerToLogsStorage := makeLogsStorageHooks()
	AddHook(hook)
	url := "https://example.com/path/to/resource"
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	SugaredWithRequest(req).Infof("success getting %s", url)

	storage := **pointerToLogsStorage
	assert.Equal(t, 1, len(storage))
	assert.Equal(t, zap.InfoLevel, storage[0].Level)
	assert.Equal(t, fmt.Sprintf("success getting %s", url), storage[0].Message)
}
