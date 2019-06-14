package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"***REMOVED***/darkroom/core/pkg/config"
	"***REMOVED***/darkroom/core/pkg/logger"
)

type Server struct {
	handler http.Handler
	hook    *LifeCycleHook
}

func NewServer(opts ...Option) *Server {
	s := Server{}
	for _, opt := range opts {
		opt(&s)
	}
	return &s
}

func (s *Server) AddLifeCycleHook(hook *LifeCycleHook) {
	s.hook = hook
}

func (s *Server) Start() {
	logger.Infof("Starting %s server", config.AppName())

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
		logger.Errorf("error while starting %s server: %s", config.AppName(), err)
	}
}

func waitForShutdown(s *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGINT,
		syscall.SIGTERM)
	_ = <-sig
	logger.Infof("%s server shutting down", config.AppName())

	err := s.Shutdown(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}
	close(sig)
	logger.Infof("%s server shutdown complete", config.AppName())
}
