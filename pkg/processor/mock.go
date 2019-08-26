package processor

import (
	"image"

	"github.com/stretchr/testify/mock"
)

type MockProcessor struct {
	mock.Mock
}

func (m *MockProcessor) Crop(img image.Image, width, height int, point CropPoint) image.Image {
	args := m.Called(img, width, height, point)
	return args.Get(0).(image.Image)
}

func (m *MockProcessor) Resize(img image.Image, width, height int) image.Image {
	args := m.Called(img, width, height)
	return args.Get(0).(image.Image)
}

func (m *MockProcessor) Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error) {
	args := m.Called(base, overlay, opacity)
	return args.Get(0).([]byte), args.Get(1).(error)
}

func (m *MockProcessor) GrayScale(img image.Image) image.Image {
	args := m.Called(img)
	return args.Get(0).(image.Image)
}

func (m *MockProcessor) Flip(img image.Image, mode string) image.Image {
	args := m.Called(img, mode)
	return args.Get(0).(image.Image)
}

func (m *MockProcessor) Rotate(img image.Image, angle float64) image.Image {
	args := m.Called(img, angle)
	return args.Get(0).(image.Image)
}

func (m *MockProcessor) Decode(data []byte) (image.Image, string, error) {
	args := m.Called(data)
	img := args.Get(0)
	ext := args.Get(1)
	if img != nil && ext != nil {
		return img.(image.Image), ext.(string), args.Error(2)
	}
	return nil, "", args.Error(2)
}

func (m *MockProcessor) Encode(img image.Image, format string, enforceFmt bool) ([]byte, error) {
	args := m.Called(img, format, enforceFmt)
	b := args.Get(0).([]byte)
	if args.Get(1) == nil {
		return b, nil
	}
	return b, args.Get(1).(error)
}

func (m *MockProcessor) FixOrientation(img image.Image, orientation int) image.Image {
	args := m.Called(img, orientation)
	return args.Get(0).(image.Image)
}

func (m *MockProcessor) Support(format string) bool {
	args := m.Called(format)
	return args.Get(0).(bool)
}
