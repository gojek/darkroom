// Package router has the default routes information
package router

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"net/http/pprof"

	"github.com/gojek/darkroom/pkg/regex"

	"github.com/gojek/darkroom/internal/handler"
	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/service"
	"github.com/gorilla/mux"
)

// NewRouter takes in handler Dependencies and returns mux.Router with default routes
// and if debug mode is enabled then it also adds pprof routes.
// It also, adds a PathPrefix to catch all route if config.DataSource().PathPrefix is set
func NewRouter(deps *service.Dependencies) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.Methods(http.MethodGet).Path("/ping").Handler(handler.Ping())

	if config.DebugModeEnabled() {
		setDebugRoutes(r)
	}
	r.Handle("/metrics", promhttp.Handler())
	// Catch all handler
	s := config.DataSource()
	if (regex.S3Matcher.MatchString(s.Kind) ||
		regex.CloudfrontMatcher.MatchString(s.Kind)) &&
		s.PathPrefix != "" {
		r.Methods(http.MethodGet).PathPrefix(s.PathPrefix).Handler(handler.ImageHandler(deps))
	} else {
		r.Methods(http.MethodGet).PathPrefix("/").Handler(handler.ImageHandler(deps))
	}

	return r
}

func setDebugRoutes(r *mux.Router) {
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.HandleFunc("/debug/pprof/goroutine", pprof.Index)
	r.HandleFunc("/debug/pprof/heap", pprof.Index)
	r.HandleFunc("/debug/pprof/threadcreate", pprof.Index)
}
