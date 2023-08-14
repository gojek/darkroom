package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSpecBuilder_Build(t *testing.T) {
	scope := "scope"
	img := []byte("imageData")
	params := map[string]string{"foo": "bar"}
	formats := []string{"image/webp", "image/apng"}
	ext := "png"
	spec := NewSpecBuilder().
		WithScope(scope).
		WithImageData(img).
		WithParams(params).
		WithFormats(formats).
		WithTargetFormat(ext).
		Build()

	assert.Equal(t, spec.Scope, scope)
	assert.Equal(t, spec.ImageData, img)
	assert.Equal(t, spec.Params, params)
	assert.Equal(t, spec.formats, formats)
	assert.Equal(t, spec.TargetFormat, ext)
}

func TestSpec_IsWebPSupported(t *testing.T) {
	f := []string{"image/webp", "image/apng"}
	spec := NewSpecBuilder().WithFormats(f).Build()
	assert.True(t, spec.IsWebPSupported())

	f = []string{"image/apng"}
	spec = NewSpecBuilder().WithFormats(f).Build()
	assert.False(t, spec.IsWebPSupported())
}

func TestSpec_Build_TargetExtensionNotValid(t *testing.T) {
	ext := "gif"
	spec := NewSpecBuilder().WithTargetFormat(ext).Build()
	assert.Empty(t, spec.TargetFormat)
}
