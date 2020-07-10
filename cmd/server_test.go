package cmd

import (
	"fmt"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	// setup
	errCh := make(chan error)
	stopCh := make(chan struct{})
	diagnosticsPort := 9999
	v := config.Viper()
	v.Set("port", diagnosticsPort)
	v.Set("source.kind", "WebFolder")
	v.Set("source.baseURL", "https://example.com/path/to/folder")
	config.Update()

	// given
	cmd := newRunCmdWithOpts(runCmdOpts{
		SetupSignalHandler: func() <-chan struct{} {
			return stopCh
		},
	})

	// when
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

	// when
	close(stopCh)

	// then
	assert.NoError(t, <-errCh)
}
