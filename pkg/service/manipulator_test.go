package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"***REMOVED***/darkroom/core/pkg/processor"
	"***REMOVED***/darkroom/core/pkg/processor/native"
	"testing"
)

func TestNewManipulator(t *testing.T) {
	m := NewManipulator(native.NewBildProcessor())
	assert.NotNil(t, m)
}

func TestManipulator_Process(t *testing.T) {
	mp := &mockProcessor{}
	m := NewManipulator(mp)
	params := make(map[string]string)

	mp.On("Crop", []byte("toCrop"), 100, 100, processor.CropCenter).Return([]byte("cropped"), nil)

	params[fit] = crop
	params[width] = "100"
	params[height] = "100"
	data, err := m.Process(context.TODO(), []byte("toCrop"), params)
	assert.Nil(t, err)
	assert.Equal(t, []byte("cropped"), data)

	mp.On("Resize", []byte("toResize"), 100, 100).Return([]byte("reSized"), nil)

	params = make(map[string]string)
	params[width] = "100"
	params[height] = "100"
	data, err = m.Process(context.TODO(), []byte("toResize"), params)
	assert.Nil(t, err)
	assert.Equal(t, []byte("reSized"), data)

	mp.On("GrayScale", []byte("toGrayScale")).Return([]byte("grayScaled"), nil)

	params = make(map[string]string)
	params[mono] = blackHexCode
	data, err = m.Process(context.TODO(), []byte("toGrayScale"), params)
	assert.Nil(t, err)
	assert.Equal(t, []byte("grayScaled"), data)
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

type mockProcessor struct {
	mock.Mock
}

func (m *mockProcessor) Crop(input []byte, width, height int, point processor.CropPoint) ([]byte, error) {
	args := m.Called(input, width, height, point)
	if args.Get(1) == nil {
		return args.Get(0).([]byte), nil
	}
	return args.Get(0).([]byte), args.Get(1).(error)
}

func (m *mockProcessor) Resize(input []byte, width, height int) ([]byte, error) {
	args := m.Called(input, width, height)
	if args.Get(1) == nil {
		return args.Get(0).([]byte), nil
	}
	return args.Get(0).([]byte), args.Get(1).(error)
}

func (m *mockProcessor) Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error) {
	args := m.Called(base, overlay, opacity)
	return args.Get(0).([]byte), args.Get(1).(error)
}

func (m *mockProcessor) GrayScale(input []byte) ([]byte, error) {
	args := m.Called(input)
	if args.Get(1) == nil {
		return args.Get(0).([]byte), nil
	}
	return args.Get(0).([]byte), args.Get(1).(error)
}
