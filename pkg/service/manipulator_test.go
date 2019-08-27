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

func TestManipulator_Process(t *testing.T) {
	mp := &mockProcessor{}
	m := NewManipulator(mp)
	params := make(map[string]string)

	input := []byte("inputData")
	decoded := &image.RGBA{Pix: []uint8{1, 2, 3, 4}}

	// Test flow for Decode error from Processor
	mp.On("Decode", mock.Anything).Return(nil, "", errors.New("decoding error"))
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())
	mp.AssertExpectations(t)

	// Create new struct for asserting expectations
	mp = &mockProcessor{}
	m = NewManipulator(mp)
	mp.On("Decode", input).Return(decoded, "png", nil)
	mp.On("Encode", decoded, "png").Return(input, nil)
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
	_, _ = m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})

	mp.On("Flip", decoded, "v").Return(decoded, nil)
	params = make(map[string]string)
	params[flip] = "v"
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("Rotate", decoded, 90.5).Return(decoded, nil)
	params = make(map[string]string)
	params[rotate] = "90.5"
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	mp.On("FixOrientation", decoded, 0).Return(decoded)
	params = make(map[string]string)
	params[auto] = compress
	_, _ = m.Process(NewSpecBuilder().WithImageData(input).WithParams(params).Build())

	// Assert all expectations once here
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

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) Crop(img image.Image, width, height int, point processor.CropPoint) image.Image {
	args := m.Called(img, width, height, point)
	return args.Get(0).(image.Image)
}

func (m *mockProcessor) Resize(img image.Image, width, height int) image.Image {
	args := m.Called(img, width, height)
	return args.Get(0).(image.Image)
}

func (m *mockProcessor) Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error) {
	args := m.Called(base, overlay, opacity)
	return args.Get(0).([]byte), args.Get(1).(error)
}

func (m *mockProcessor) GrayScale(img image.Image) image.Image {
	args := m.Called(img)
	return args.Get(0).(image.Image)
}

func (m *mockProcessor) Blur(img image.Image, radius float64) image.Image {
	args := m.Called(img, radius)
	return args.Get(0).(image.Image)
}

func (m *mockProcessor) Flip(img image.Image, mode string) image.Image {
	args := m.Called(img, mode)
	return args.Get(0).(image.Image)
}

func (m *mockProcessor) Rotate(img image.Image, angle float64) image.Image {
	args := m.Called(img, angle)
	return args.Get(0).(image.Image)
}

func (m *mockProcessor) Decode(data []byte) (image.Image, string, error) {
	args := m.Called(data)
	img := args.Get(0)
	ext := args.Get(1)
	if img != nil && ext != nil {
		return img.(image.Image), ext.(string), args.Error(2)
	}
	return nil, "", args.Error(2)
}

func (m *mockProcessor) Encode(img image.Image, format string) ([]byte, error) {
	args := m.Called(img, format)
	b := args.Get(0).([]byte)
	if args.Get(1) == nil {
		return b, nil
	}
	return b, args.Get(1).(error)
}

func (m *mockProcessor) FixOrientation(img image.Image, orientation int) image.Image {
	args := m.Called(img, orientation)
	return args.Get(0).(image.Image)
}
