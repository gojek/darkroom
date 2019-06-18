package native

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"***REMOVED***/darkroom/core/pkg/processor"
	"testing"
)

const (
	actualWidth  = 1000
	actualHeight = 500
)

func Test_getResizeWidthAndHeight_ZeroHeight(t *testing.T) {
	w, h := getResizeWidthAndHeight(2000, 0, actualWidth, actualHeight)
	assert.Equal(t, 2000, w)
	assert.Equal(t, 1000, h)
}

func Test_getResizeWidthAndHeight_ZeroWidth(t *testing.T) {
	w, h := getResizeWidthAndHeight(0, 750, actualWidth, actualHeight)
	assert.Equal(t, 1500, w)
	assert.Equal(t, 750, h)
}

func Test_getResizeWidthAndHeight(t *testing.T) {
	w, h := getResizeWidthAndHeight(1090, 470, actualWidth, actualHeight)
	assert.Equal(t, 940, w)
	assert.Equal(t, 470, h)
	w, h = getResizeWidthAndHeight(200, 300, actualWidth, actualHeight)
	assert.Equal(t, 200, w)
	assert.Equal(t, 100, h)
}

func TestGetResizeWidthAndHeightForFitCrop(t *testing.T) {
	w, h := getResizeWidthAndHeightForCrop(800, 400, actualWidth, actualHeight)
	assert.Equal(t, 800, w)
	assert.Equal(t, 400, h)

	w, h = getResizeWidthAndHeightForCrop(200, 0, actualWidth, actualHeight)
	assert.Equal(t, 200, w)
	assert.Equal(t, 100, h)

	w, h = getResizeWidthAndHeightForCrop(0, 300, actualWidth, actualHeight)
	assert.Equal(t, 600, w)
	assert.Equal(t, 300, h)

	w, h = getResizeWidthAndHeightForCrop(200, 300, actualWidth, actualHeight)
	assert.Equal(t, 600, w)
	assert.Equal(t, 300, h)
}

func TestGetStartingPointForCrop(t *testing.T) {
	//center
	x, y := getStartingPointForCrop(500, 500, 300, 500, processor.CropCenter)
	assert.Equal(t, 100, x)
	assert.Equal(t, 0, y)

	//top
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropTop)
	assert.Equal(t, 100, x)
	assert.Equal(t, 0, y)

	//topLeft
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropTopLeft)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	//topRight
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropTopRight)
	assert.Equal(t, 200, x)
	assert.Equal(t, 0, y)

	//left
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropLeft)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	//right
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropRight)
	assert.Equal(t, 200, x)
	assert.Equal(t, 0, y)

	//bottom
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropBottom)
	assert.Equal(t, 100, x)
	assert.Equal(t, 0, y)

	//bottomLeft
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropBottomLeft)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	//bottomRight
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.CropBottomRight)
	assert.Equal(t, 200, x)
	assert.Equal(t, 0, y)
}

func Test_isOpaqueWithoutOpaqueMethod(t *testing.T) {
	im := &mockImage{opaque: false}
	val := isOpaque(im)
	assert.False(t, val)

	im = &mockImage{opaque: true}
	val = isOpaque(im)
	assert.True(t, val)
}

type mockImage struct {
	opaque bool
}

func (im *mockImage) ColorModel() color.Model {
	panic("implement me")
}

func (im *mockImage) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{X: 5, Y: 5},
		Max: image.Point{X: 10, Y: 10},
	}
}

func (im *mockImage) At(x, y int) color.Color {
	return &mockColor{opaque: im.opaque}
}

type mockColor struct {
	opaque bool
}

func (m *mockColor) RGBA() (r, g, b, a uint32) {
	if m.opaque {
		return 0x0fff, 0xf0ff, 0xff0f, 0xffff
	}
	return 0x0fff, 0xf0ff, 0xff0f, 0xfff0
}
