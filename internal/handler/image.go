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
	ContentLengthHeader = "Content-Length"
	CacheControlHeader  = "Cache-Control"
	StorageGetErrorKey  = "handler.storage.get.error"
	ProcessorErrorKey   = "handler.processor.error"
)

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
