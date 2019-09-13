package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gojek/darkroom/pkg/metrics"
	"github.com/gojek/darkroom/pkg/service"
)

const (
	// ContentLengthHeader is the response header key used to set content length
	ContentLengthHeader = "Content-Length"
	// CacheControlHeader is the response header key used to set cache control
	CacheControlHeader = "Cache-Control"
	// VaryHeader is the response header key used to indicate the CDN that the response should depend on client's accept header
	// Ref: https://tools.ietf.org/html/rfc7231#section-7.1.4
	VaryHeader = "Vary"
	// StorageGetErrorKey is the key used while pushing metrics update to statsd
	StorageGetErrorKey = "handler.storage.get.error"
	// ProcessorErrorKey is the key used while pushing metrics update to statsd
	ProcessorErrorKey = "handler.processor.error"
)

// ImageHandler is responsible for fetching the path from the storage backend and processing it if required
func ImageHandler(deps *service.Dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logger.SugaredWithRequest(r)
		res := deps.Storage.Get(r.Context(), r.URL.Path)
		if res.Error() != nil {
			l.Errorf("error from Storage.Get: %s", res.Error())
			metrics.Update(metrics.UpdateOption{Name: StorageGetErrorKey, Type: metrics.Count})
			w.WriteHeader(res.Status())
			return
		}
		var data []byte
		var err error
		data = res.Data()

		params := getParams(r.URL.Query())
		if len(*params) > 0 {
			data, err = deps.Manipulator.Process(service.NewSpecBuilder().WithImageData(data).WithParams(*params).Build())
			if err != nil {
				l.Errorf("error from Manipulator.Process: %s", err)
				metrics.Update(metrics.UpdateOption{Name: ProcessorErrorKey, Type: metrics.Count})
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		cl, _ := w.Write(data)
		w.Header().Set(ContentLengthHeader, fmt.Sprintf("%d", cl))
		w.Header().Set(CacheControlHeader, fmt.Sprintf("public,max-age=%d", config.CacheTime()))
		// Ref to Google CDN we support: https://cloud.google.com/cdn/docs/caching#cacheability
		w.Header().Set(VaryHeader, "Accept")
	}
}

func getParams(values url.Values) *map[string]string {
	params := make(map[string]string)
	for _, param := range config.DefaultParams() {
		if strings.Contains(param, "=") {
			p := strings.Split(param, "=")
			params[p[0]] = p[1]
		}
	}
	for v := range values {
		if len(values.Get(v)) != 0 {
			params[v] = values.Get(v)
		}
	}
	return &params
}
