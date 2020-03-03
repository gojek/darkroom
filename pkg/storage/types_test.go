package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewResponse(t *testing.T) {
	err := errors.New("randomError")
	metadata := &ResponseMetadata{
		AcceptRanges:  "bytes",
		ContentLength: "101",
		ContentRange:  "bytes 100-200/247103",
		ContentType:   "image/png",
		ETag:          "32705ce195789d7bf07f3d44783c2988",
		LastModified:  "Wed, 21 Oct 2015 07:28:00 GMT ",
	}
	r := NewResponse([]byte("randomBytes"), http.StatusBadRequest, err).WithMetadata(metadata)

	assert.Equal(t, []byte("randomBytes"), r.Data())
	assert.Equal(t, http.StatusBadRequest, r.Status())
	assert.Equal(t, err, r.Error())
	assert.Equal(t, metadata, r.Metadata())
}
