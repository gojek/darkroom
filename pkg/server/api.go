package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/logger"
)

// Server struct wraps a http.Handler and the LifeCycleHook
type Server struct {
	handler http.Handler
	hook    *LifeCycleHook
}

// NewServer returns a new Server configurable with Options
func NewServer(opts ...Option) *Server {
	s := Server{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

// AddLifeCycleHook sets the passed LifeCycleHook to the Server struct
func (s *Server) AddLifeCycleHook(hook *LifeCycleHook) {
	s.hook = hook
}

// Start is used to start a http.Server and wait for a kill signal to gracefully shutdown the server
func (s *Server) Start() {
	logger.Info("Starting darkroom server")

	if s.hook != nil {
		s.hook.initFunc()
		defer s.hook.deferFunc()
	}

	portInfo := fmt.Sprintf(":%d", config.Port())
	server := &http.Server{Addr: portInfo, Handler: s.handler}

	go listenServer(server)
	waitForShutdown(server)
}

func listenServer(s *http.Server) {
	err := s.ListenAndServe()
	if err != http.ErrServerClosed && err != nil {
		logger.Errorf("error while starting darkroom server: %s", err)
	}
}

func waitForShutdown(s *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	logger.Info("darkroom server shutting down")

	err := s.Shutdown(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}
	close(sig)
	logger.Info("darkroom server shutdown complete")
}
