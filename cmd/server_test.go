package cmd

import (
	"fmt"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	errCh := make(chan error)
	stopCh := make(chan struct{})
	diagnosticsPort := 9999
	v := config.Viper()
	v.Set("source.kind", "WebFolder")
	v.Set("source.baseURL", "https://example.com/path/to/folder")
	v.Set("metrics.system", "prometheus")
	config.Update()

	cmd := newRunCmdWithOpts(runCmdOpts{
		SetupSignalHandler: func() <-chan struct{} {
			return stopCh
		},
		registry: prometheus.NewRegistry(),
	})
	cmd.SetArgs([]string{"-p", fmt.Sprintf("%d", diagnosticsPort)})

	go func() {
		defer close(errCh)
		errCh <- cmd.Execute()
	}()

	assert.True(t, assert.Eventually(t, func() bool {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/ping", diagnosticsPort))
		if err != nil {
			return false
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 100*time.Millisecond), "failed to run server")

	close(stopCh)

	assert.NoError(t, <-errCh)
}

func TestRunServerWithInvalidPort(t *testing.T) {
	errCh := make(chan error)
	stopCh := make(chan struct{})
	v := config.Viper()
	v.Set("source.kind", "WebFolder")
	v.Set("source.baseURL", "https://example.com/path/to/folder")
	v.Set("metrics.system", "prometheus")
	config.Update()

	cmd := newRunCmdWithOpts(runCmdOpts{
		SetupSignalHandler: func() <-chan struct{} {
			return stopCh
		},
		registry: prometheus.NewRegistry(),
	})
	cmd.SetArgs([]string{"-p", fmt.Sprintf("%d", -9000)})

	go func() {
		defer close(errCh)
		errCh <- cmd.Execute()
	}()

	assert.Error(t, <-errCh)
}

func TestRunServerWithInvalidDependencies(t *testing.T) {
	errCh := make(chan error)
	stopCh := make(chan struct{})
	v := config.Viper()
	v.Set("source.kind", "")
	v.Set("source.baseURL", "")
	v.Set("metrics.system", "")
	config.Update()

	cmd := newRunCmdWithOpts(runCmdOpts{
		SetupSignalHandler: func() <-chan struct{} {
			return stopCh
		},
		registry: prometheus.NewRegistry(),
	})

	go func() {
		defer close(errCh)
		errCh <- cmd.Execute()
	}()

	assert.EqualError(t, <-errCh, "handler dependencies are not valid")
}
