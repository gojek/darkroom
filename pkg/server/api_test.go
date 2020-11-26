package server

import (
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	s := NewServer(Options{
		Handler:       mux.NewRouter(),
		Port:          3000,
		LifeCycleHook: nil,
	})
	assert.NotNil(t, s)
}

func TestLifeCycleHook(t *testing.T) {
	init, deferred := false, false
	s := NewServer(Options{
		Handler: mux.NewRouter(),
		Port:    3000,
		LifeCycleHook: &LifeCycleHook{
			initFunc:  func() { init = true },
			deferFunc: func() { deferred = true },
		},
	})
	assert.NotNil(t, s)
	stopCh := make(chan struct{})
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		err := s.Start(stopCh)
		errCh <- err
	}()

	// wait for above goroutine to run
	time.Sleep(1 * time.Second)
	assert.True(t, init)
	assert.False(t, deferred)
	close(stopCh)

	<-errCh
	assert.True(t, deferred)
}
