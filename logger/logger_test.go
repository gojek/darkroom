package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestLogger(t *testing.T) {
	hook := test.NewLocal(getLogger())
	AddHook(hook)

	cases := []struct {
		message  string
		level    logrus.Level
		callFunc func(args ...interface{})
	}{
		{
			message:  "warning message",
			level:    logrus.WarnLevel,
			callFunc: Warn,
		},
		{
			message:  "warning message in new line",
			level:    logrus.WarnLevel,
			callFunc: Warnln,
		},
		{
			message:  "debug message",
			level:    logrus.DebugLevel,
			callFunc: Debug,
		},
		{
			message:  "debug message in new line",
			level:    logrus.DebugLevel,
			callFunc: Debugln,
		},
		{
			message:  "error message",
			level:    logrus.ErrorLevel,
			callFunc: Error,
		},
		{
			message:  "debug message in new line",
			level:    logrus.ErrorLevel,
			callFunc: Errorln,
		},
		{
			message:  "info message",
			level:    logrus.InfoLevel,
			callFunc: Info,
		},
		{
			message:  "info message in new line",
			level:    logrus.InfoLevel,
			callFunc: Infoln,
		},
	}

	for i, c := range cases {
		c.callFunc(c.message)
		assert.Equal(t, (i+1)*2, len(hook.Entries))
		assert.Equal(t, c.level, hook.LastEntry().Level)
		assert.Equal(t, c.message, hook.LastEntry().Message)
	}
}

func TestLoggerFormatter(t *testing.T) {
	hook := test.NewLocal(getLogger())
	AddHook(hook)

	cases := []struct {
		format   string
		level    logrus.Level
		args     []string
		callFunc func(format string, args ...interface{})
	}{
		{
			format:   "warning message: %s",
			level:    logrus.WarnLevel,
			args:     []string{"WARNING"},
			callFunc: Warnf,
		},
		{
			format:   "debug message: %s",
			level:    logrus.DebugLevel,
			args:     []string{"DEBUG"},
			callFunc: Debugf,
		},
		{
			format:   "error message: %s",
			level:    logrus.ErrorLevel,
			args:     []string{"ERROR"},
			callFunc: Errorf,
		},
		{
			format:   "info message: %s",
			level:    logrus.InfoLevel,
			args:     []string{"INFO"},
			callFunc: Infof,
		},
	}

	for i, c := range cases {
		c.callFunc(c.format, c.args)
		assert.Equal(t, (i+1)*2, len(hook.Entries))
		assert.Equal(t, c.level, hook.LastEntry().Level)
	}
}

func TestLoggerWithFields(t *testing.T) {
	hook := test.NewLocal(getLogger())
	AddHook(hook)

	msg := "WithFields"

	fields := make(map[string]interface{})
	fields["key"] = "value"
	WithFields(fields).Warn(msg)

	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)
}

func TestLoggerWithField(t *testing.T) {
	hook := test.NewLocal(getLogger())
	AddHook(hook)

	msg := "WithField"

	WithField("key", "value").Warn(msg)

	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)
}

func TestLoggerWithRequest(t *testing.T) {
	hook := test.NewLocal(getLogger())
	AddHook(hook)
	msg := "success"
	req, _ := http.NewRequest(http.MethodGet, "https://example.com/path/to/resource", nil)

	WithRequest(req).Info(msg)

	assert.Equal(t, 2, len(hook.Entries))
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)
}
