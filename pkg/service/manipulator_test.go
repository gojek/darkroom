package service

import (
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

	mp.On("Decode", input).Return(decoded, "png", nil)
	mp.On("Encode", decoded, "png").Return(input, nil)
	mp.On("Crop", decoded, 100, 100, processor.CropCenter).Return(decoded, nil)

	params[fit] = crop
	params[width] = "100"
	params[height] = "100"
	data, err := m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})
	assert.Nil(t, err)
	assert.Equal(t, input, data)

	mp.On("Resize", decoded, 100, 100).Return(decoded, nil)

	params = make(map[string]string)
	params[width] = "100"
	params[height] = "100"
	data, err = m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})
	assert.Nil(t, err)
	assert.Equal(t, input, data)

	mp.On("GrayScale", decoded).Return(decoded, nil)

	params = make(map[string]string)
	params[mono] = blackHexCode
	data, err = m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})
	assert.Nil(t, err)
	assert.Equal(t, input, data)

	mp.On("Flip", decoded, "v").Return(decoded, nil)

	params = make(map[string]string)
	params[flip] = "v"
	data, err = m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})
	assert.Nil(t, err)
	assert.Equal(t, input, data)

	mp.On("Rotate", decoded, 90.5).Return(decoded, nil)

	params = make(map[string]string)
	params[rotate] = "90.5"
	data, err = m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})
	assert.Nil(t, err)
	assert.Equal(t, input, data)

	mp.On("FixOrientation", decoded, mock.Anything).Return(decoded)
	params = make(map[string]string)
	params[auto] = compress
	data, err = m.Process(ProcessSpec{
		ImageData: input,
		Params:    params,
	})
	assert.Nil(t, err)
	assert.True(t, mp.AssertCalled(t, "FixOrientation", decoded, 0))
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

	updateOption := trackDuration(cropDurationKey, time.Now(), ProcessSpec{
		ImageData: imageData,
	})
	assert.Equal(t, fmt.Sprintf("%s.%s.%s", cropDurationKey, "<=128KB", "png"), updateOption.Name)

	updateOption = trackDuration(cropDurationKey, time.Now(), ProcessSpec{
		ImageData: make([]byte, 10, 10),
	})
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
	img := args.Get(0).(image.Image)
	ext := args.Get(1).(string)
	if args.Get(2) == nil {
		return img, ext, nil
	}
	return img, ext, args.Get(2).(error)
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
