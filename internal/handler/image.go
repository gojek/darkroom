package handler

import (
	"fmt"
	"github.com/gojek/darkroom/pkg/metrics"
	"net/http"

	"github.com/gojek/darkroom/pkg/config"
	"github.com/gojek/darkroom/pkg/logger"
	"github.com/gojek/darkroom/pkg/service"
)

const (
	// ContentLengthHeader is the response header key used to set content length
	ContentLengthHeader = "Content-Length"
	// CacheControlHeader is the response header key used to set cache controll
	CacheControlHeader  = "Cache-Control"
	// StorageGetErrorKey is the key used while pushing metrics update to statsd
	StorageGetErrorKey  = "handler.storage.get.error"
	// ProcessorErrorKey is the key used while pushing metrics update to statsd
	ProcessorErrorKey   = "handler.processor.error"
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

		params := make(map[string]string)
		values := r.URL.Query()
		if len(values) > 0 {
			for v := range values {
				if len(values.Get(v)) != 0 {
					params[v] = values.Get(v)
				}
			}
			data, err = deps.Manipulator.Process(service.ProcessSpec{
				ImageData: data,
				Params:    params,
			})
			if err != nil {
				l.Errorf("error from Manipulator.Process: %s", err)
				metrics.Update(metrics.UpdateOption{Name: ProcessorErrorKey, Type: metrics.Count})
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}
		}

		cl, _ := w.Write([]byte(data))
		w.Header().Set(ContentLengthHeader, string(cl))
		w.Header().Set(CacheControlHeader, fmt.Sprintf("public,max-age=%d", config.CacheTime()))
	}
}
