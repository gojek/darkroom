package native

import (
	"github.com/gojek/darkroom/pkg/config"
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
	x, y := getStartingPointForCrop(500, 500, 300, 500, processor.PointCenter)
	assert.Equal(t, 100, x)
	assert.Equal(t, 0, y)

	//top
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointTop)
	assert.Equal(t, 100, x)
	assert.Equal(t, 0, y)

	//topLeft
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointTopLeft)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	//topRight
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointTopRight)
	assert.Equal(t, 200, x)
	assert.Equal(t, 0, y)

	//left
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointLeft)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	//right
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointRight)
	assert.Equal(t, 200, x)
	assert.Equal(t, 0, y)

	//bottom
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointBottom)
	assert.Equal(t, 100, x)
	assert.Equal(t, 0, y)

	//bottomLeft
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointBottomLeft)
	assert.Equal(t, 0, x)
	assert.Equal(t, 0, y)

	//bottomRight
	x, y = getStartingPointForCrop(500, 500, 300, 500, processor.PointBottomRight)
	assert.Equal(t, 200, x)
	assert.Equal(t, 0, y)
}

func Test_isOpaqueWithFastOpaqueMethod(t *testing.T) {
	r := image.Rect(0, 0, 640, 480)
	gray, gray16, cmyk := image.NewGray(r), image.NewGray16(r), image.NewCMYK(r)
	assert.True(t, isOpaque(gray))
	assert.True(t, isOpaque(gray16))
	assert.True(t, isOpaque(cmyk))
}

func Test_isOpaqueWithoutFastOpaqueMethodShouldReturnTrue(t *testing.T) {
	isOpaqueShouldReturnTrue := func() {
		img := NewMockImage(image.Rect(0, 0, 640, 480))
		draw.Draw(img, img.Bounds(), image.Opaque, image.ZP, draw.Src)
		val := isOpaque(img)
		assert.True(t, val)
	}
	v := config.Viper()
	v.Set("enableConcurrentOpacityChecking", true)
	config.Update()
	isOpaqueShouldReturnTrue()
	v.Set("enableConcurrentOpacityChecking", false)
	isOpaqueShouldReturnTrue()
}

func Test_isOpaqueWithoutFastOpaqueMethodShouldReturnFalse(t *testing.T) {
	isOpaqueShouldReturnFalse := func() {
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
	v := config.Viper()
	v.Set("enableConcurrentOpacityChecking", true)
	config.Update()
	isOpaqueShouldReturnFalse()
	v.Set("enableConcurrentOpacityChecking", false)
	isOpaqueShouldReturnFalse()
}

type MockImage struct {
	rect   image.Rectangle
	points [][]color.Color
}

func NewMockImage(rect image.Rectangle) *MockImage {
	mockImg := &MockImage{
		rect: rect,
	}
	points := make([][]color.Color, rect.Dy())
	for i := range points {
		points[i] = make([]color.Color, rect.Dx())
	}
	mockImg.points = points
	return mockImg
}

func (im *MockImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (im *MockImage) Bounds() image.Rectangle {
	return im.rect
}

func (im *MockImage) At(x, y int) color.Color {
	return im.points[y][x]
}

func (im *MockImage) Set(x, y int, c color.Color) {
	im.points[y][x] = c
}
