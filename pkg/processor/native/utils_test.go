package native

import (
	"github.com/gojek/darkroom/pkg/processor"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/draw"
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

func Test_isOpaqueWithoutOpaqueMethodShouldReturnTrue(t *testing.T) {
	img := NewMockImage(image.Rect(0, 0, 640, 480))
	draw.Draw(img, img.Bounds(), image.Opaque, image.ZP, draw.Src)
	val := isOpaque(img)
	assert.True(t, val)
}

func Test_isOpaqueWithoutOpaqueMethodShouldReturnFalse(t *testing.T) {
	w, h := 640, 480
	img := NewMockImage(image.Rect(0, 0, w, h))
	draw.Draw(img, img.Bounds(), image.Opaque, image.ZP, draw.Src)

	cases := []struct {
		x, y int
	}{
		{x: 0, y: 0},
		{x: w / 2, y: h / 2},
		{x: w - 1, y: h - 1},
	}
	for _, c := range cases {
		// Flip only 1 bit to be transparent for each test case
		x, y := c.x, c.y
		img.Set(x, y, image.Transparent.C)
		val := isOpaque(img)
		assert.False(t, val)
		img.Set(x, y, image.Opaque.C)
	}
}

type mockImage struct {
	rect   image.Rectangle
	points [][]color.Color
}

func NewMockImage(rect image.Rectangle) *mockImage {
	mockImg := &mockImage{
		rect: rect,
	}
	points := make([][]color.Color, rect.Dy())
	for i := range points {
		points[i] = make([]color.Color, rect.Dx())
	}
	mockImg.points = points
	return mockImg
}

func (im *mockImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (im *mockImage) Bounds() image.Rectangle {
	return im.rect
}

func (im *mockImage) At(x, y int) color.Color {
	return im.points[y][x]
}

func (im *mockImage) Set(x, y int, c color.Color) {
	im.points[y][x] = c
}
