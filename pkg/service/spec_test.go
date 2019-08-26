package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecBuilder_Build(t *testing.T) {
	scope := "scope"
	img := []byte("imageData")
	params := make(map[string]string)
	params["foo"] = "bar"
	formats := []string{"image/webp", "image/apng"}

	spec := NewSpecBuilder().
		WithScope(scope).
		WithImageData(img).
		WithParams(params).
		WithFormats(formats).
		Build()

	assert.Equal(t, spec.Scope, scope)
	assert.Equal(t, spec.ImageData, img)
	assert.Equal(t, spec.Params, params)
	assert.Equal(t, spec.Formats, formats)
}
