package native

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"image"
	"io/ioutil"
	"***REMOVED***/darkroom/core/pkg/processor"
	"testing"
)

type BildProcessorSuite struct {
	suite.Suite
	srcData       []byte
	watermarkData []byte
	badData       []byte
	processor     processor.Processor
}

func (s *BildProcessorSuite) SetupSuite() {
	s.processor = NewBildProcessor()
	s.srcData, _ = ioutil.ReadFile("test.png")
	s.watermarkData, _ = ioutil.ReadFile("overlay.png")
	s.badData = []byte("badImage.ext")
}

func TestNewBildProcessor(t *testing.T) {
	suite.Run(t, new(BildProcessorSuite))
}

func (s *BildProcessorSuite) TestBildProcessor_Resize() {
	output, err := s.processor.Resize(s.srcData, 500, 500)

	assert.NotNil(s.T(), output)
	assert.Nil(s.T(), err)

	img, _, _ := image.Decode(bytes.NewReader(output))
	assert.Equal(s.T(), 500, img.Bounds().Dx())
	assert.Equal(s.T(), 375, img.Bounds().Dy())
}

func (s *BildProcessorSuite) TestBildProcessor_Crop() {
	output, err := s.processor.Crop(s.srcData, 500, 500, processor.CropCenter)

	assert.NotNil(s.T(), output)
	assert.Nil(s.T(), err)

	img, _, _ := image.Decode(bytes.NewReader(output))
	assert.Equal(s.T(), 500, img.Bounds().Dx())
	assert.Equal(s.T(), 500, img.Bounds().Dy())
}

func (s *BildProcessorSuite) TestBildProcessor_Grayscale() {
	output, err := s.processor.GrayScale(s.srcData)

	assert.NotNil(s.T(), output)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), s.srcData, output)
}

func (s *BildProcessorSuite) TestBildProcessor_Watermark() {
	output, err := s.processor.Watermark(s.srcData, s.watermarkData, 200)

	assert.NotNil(s.T(), output)
	assert.Nil(s.T(), err)

	assert.NotEqual(s.T(), s.srcData, output)
}

func (s *BildProcessorSuite) TestBildProcessorWithBadInput() {
	output, err := s.processor.Crop(s.badData, 0, 0, processor.CropCenter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)

	output, err = s.processor.Resize(s.badData, 0, 0)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)

	output, err = s.processor.GrayScale(s.badData)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)

	output, err = s.processor.Watermark(s.badData, s.watermarkData, 255)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)

	output, err = s.processor.Watermark(s.srcData, s.badData, 255)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)
}
