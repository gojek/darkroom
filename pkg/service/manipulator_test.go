// TODO: Change test case to use suite
package service

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"testing"
	"time"

	"github.com/gojek/darkroom/pkg/processor"
	"github.com/gojek/darkroom/pkg/processor/native"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewManipulator(t *testing.T) {
	m := NewManipulator(native.NewBildProcessor())
	assert.NotNil(t, m)
}

// Integration test to verify the flow of WebP image is requested without having support of WebP on client's side
func TestManipulator_Process_ReturnsImageAsPNGIfCallerDoesNOTSupportWebP(t *testing.T) {
	// Use real processor to ensure that right encoder is being used
	p := native.NewBildProcessor()
	m := NewManipulator(p)

	img, _ := ioutil.ReadFile("../processor/native/_testdata/test.webp")
	expectedImg, _ := ioutil.ReadFile("../processor/native/_testdata/test_webp_to_png.png")

	s := NewSpecBuilder().
		WithImageData(img).
		WithParams(map[string]string{auto: format}).
		Build()
	img, err := m.Process(s)
	assert.Nil(t, err)
	assert.Equal(t, expectedImg, img)
}

// Integration test to verify the flow of PNG image is requested with having support of WebP on client's side
func TestManipulator_Process_ReturnsImageAsWebPIfCallerSupportsWebP(t *testing.T) {
	// Use real processor to ensure that right encoder is being used
	p := native.NewBildProcessor()
	m := NewManipulator(p)

	img, _ := ioutil.ReadFile("../processor/native/_testdata/test.png")
	expectedImg, _ := ioutil.ReadFile("../processor/native/_testdata/test_png_to_webp.webp")

	s := NewSpecBuilder().
		WithImageData(img).
		WithParams(map[string]string{auto: format}).
		WithFormats([]string{"image/webp"}).
		Build()
	img, err := m.Process(s)
	assert.Nil(t, err)
	assert.Equal(t, expectedImg, img)
}

func TestManipulator_Process(t *testing.T) {
	// TODO: Refactor this test case to be more modular
	mp := &processor.MockProcessor{}
	m := NewManipulator(mp)
	params := make(map[string]string)

	input := []byte("inputData")
	decoded := &image.RGBA{Pix: []uint8{1, 2, 3, 4}}

	// Test flow for Decode error from Processor
	mp.On("Decode", mock.Anything).Return(nil, "", errors.New("decoding error"))
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())
	mp.AssertExpectations(t)

	// Create new struct for asserting expectations
	mp = &processor.MockProcessor{}
	m = NewManipulator(mp)
	mp.On("Support", mock.Anything).Return(false)
	mp.On("Decode", input).Return(decoded, "png", nil)
	mp.On("Encode", decoded, "png", false).Return(input, nil)
	mp.On("Crop", decoded, 100, 100, processor.CropCenter).Return(decoded, nil)
	params[fit] = crop
	params[width] = "100"
	params[height] = "100"
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("Resize", decoded, 100, 100).Return(decoded, nil)
	params = make(map[string]string)
	params[width] = "100"
	params[height] = "100"
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("GrayScale", decoded).Return(decoded, nil)
	params = make(map[string]string)
	params[mono] = blackHexCode
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("Blur", decoded, 60.0).Return(decoded, nil)
	params = make(map[string]string)
	params[blur] = "60"
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("Flip", decoded, "v").Return(decoded, nil)
	params = make(map[string]string)
	params[flip] = "v"
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("Rotate", decoded, 90.5).Return(decoded, nil)
	params = map[string]string{rotate: "90.5"}
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("FixOrientation", decoded, 0).Return(decoded)
	params = map[string]string{auto: compress}
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("Decode", input).Return(decoded, processor.FormatWebP, nil)
	params = map[string]string{auto: format}
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	// Assert all expectations once here
	mp.AssertExpectations(t)
}

func TestManipulator_Process_GivenValidFmQueryParameterShouldEncodeToCustomFormat(t *testing.T) {
	originalFmt, customFmt := "png", "jpeg"
	assert.NotEqual(t, originalFmt, customFmt)

	input, expected := []byte("input"), []byte("output")
	img := &image.RGBA{}
	mp := &processor.MockProcessor{}
	mp.On("Decode", input).Return(img, originalFmt, nil)
	mp.On("Support", customFmt).Return(true)
	// enforceFmt should be true when fm query parameter is present and valid
	mp.On("Encode", img, customFmt, true).Return(expected, nil)

	m := NewManipulator(mp)
	params := map[string]string{fm: customFmt}
	actual, err := m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
	mp.AssertExpectations(t)
}

func TestManipulator_Process_GivenInvalidFmQueryParameterShouldEncodeToOriginalFormat(t *testing.T) {
	originalFmt, customFmt := "png", "unknown"

	input, expected := []byte("input"), []byte("output")
	img := &image.RGBA{}
	mp := &processor.MockProcessor{}
	mp.On("Decode", input).Return(img, originalFmt, nil)
	mp.On("Support", customFmt).Return(false)
	mp.On("Encode", img, originalFmt, false).Return(expected, nil)

	m := NewManipulator(mp)
	params := map[string]string{fm: customFmt}
	actual, err := m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
	mp.AssertExpectations(t)
}

func TestGetCropPoint(t *testing.T) {
	assert.Equal(t, processor.CropCenter, GetCropPoint(""))
	assert.Equal(t, processor.CropTop, GetCropPoint("top"))
	assert.Equal(t, processor.CropTopLeft, GetCropPoint("top,left"))
	assert.Equal(t, processor.CropTopRight, GetCropPoint("top,right"))
	assert.Equal(t, processor.CropLeft, GetCropPoint("left"))
	assert.Equal(t, processor.CropRight, GetCropPoint("right"))
	assert.Equal(t, processor.CropBottom, GetCropPoint("bottom"))
	assert.Equal(t, processor.CropBottomLeft, GetCropPoint("bottom,left"))
	assert.Equal(t, processor.CropBottomRight, GetCropPoint("bottom,right"))
	assert.Equal(t, processor.CropCenter, GetCropPoint("random"))
}

func TestCleanInt(t *testing.T) {
	assert.Equal(t, 999, CleanInt("999"))
	assert.Equal(t, 23, CleanInt("23"))
	assert.Equal(t, 0, CleanInt("10000")) // Max value at 9999
	assert.Equal(t, 9999, CleanInt("9999"))
	assert.Equal(t, 0, CleanInt("0"))
	assert.Equal(t, 0, CleanInt("garbage"))
	assert.Equal(t, 0, CleanInt("-234"))
}

func Test_trackDuration(t *testing.T) {
	imageData, err := ioutil.ReadFile("../processor/native/_testdata/test.png")
	if err != nil {
		panic(err)
	}

	updateOption := trackDuration(cropDurationKey, time.Now(), NewSpecBuilder().WithImageData(imageData).Build())
	assert.Equal(t, fmt.Sprintf("%s.%s.%s", cropDurationKey, "<=128KB", "png"), updateOption.Name)

	updateOption = trackDuration(cropDurationKey, time.Now(), NewSpecBuilder().WithImageData(make([]byte, 10, 10)).Build())
	assert.Equal(t, fmt.Sprintf("%s.%s.%s", cropDurationKey, "<=128KB", "octet-stream"), updateOption.Name)
}
