package server

import (
	"context"
	"fmt"
	"github.com/gojek/darkroom/pkg/logger"
	"net/http"
)

// Options represents the Server options
type Options struct {
	Handler       http.Handler
	Port          int
	LifeCycleHook *LifeCycleHook
}

// Server struct wraps a http.Handler and the LifeCycleHook
type Server struct {
	server *http.Server
	hook   *LifeCycleHook
}

// NewServer returns a new Server configurable with Options
func NewServer(options Options) *Server {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", options.Port),
		Handler: options.Handler,
	}
	s := Server{
		server: srv,
		hook:   options.LifeCycleHook,
	}
	return &s
}

// Start is used to start a http.Server and wait for a kill signal to gracefully shutdown the server
func (s *Server) Start(stop <-chan struct{}) error {
	if s.hook != nil {
		s.hook.initFunc()
		defer s.hook.deferFunc()
	}

	errChan := make(chan error)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			switch err {
			case http.ErrServerClosed:
				return
			default:
				logger.Errorf("could not start server: %s", err)
				errChan <- err
			}
		}
	}()
	logger.Infof("Starting darkroom server at 0.0.0.0%s", s.server.Addr)

	select {
	case <-stop:
		logger.Info("Shutting down server")
		return s.server.Shutdown(context.Background())
	case err := <-errChan:
		return err
	}
}
