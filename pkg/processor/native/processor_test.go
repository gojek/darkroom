package native

import (
	"image"
	"io/ioutil"
	"testing"

	"github.com/gojek/darkroom/pkg/processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BildProcessorSuite struct {
	suite.Suite
	srcData       []byte
	srcImage      image.Image
	watermarkData []byte
	badData       []byte
	badImage      image.Image
	processor     processor.Processor
}

func (s *BildProcessorSuite) SetupSuite() {
	s.processor = NewBildProcessor()
	s.srcData, _ = ioutil.ReadFile("_testdata/test.png")
	s.srcImage, _, _ = s.processor.Decode(s.srcData)
	s.watermarkData, _ = ioutil.ReadFile("_testdata/overlay.png")
	s.badData = []byte("badImage.ext")
}

func TestNewBildProcessor(t *testing.T) {
	suite.Run(t, new(BildProcessorSuite))
}

func (s *BildProcessorSuite) TestBildProcessor_Resize() {
	out := s.processor.Resize(s.srcImage, 500, 500)

	assert.NotNil(s.T(), out)
	assert.Equal(s.T(), 500, out.Bounds().Dx())
	assert.Equal(s.T(), 375, out.Bounds().Dy())
}

func (s *BildProcessorSuite) TestBildProcessor_Crop() {
	out := s.processor.Crop(s.srcImage, 500, 500, processor.CropCenter)

	assert.NotNil(s.T(), out)

	assert.Equal(s.T(), 500, out.Bounds().Dx())
	assert.Equal(s.T(), 500, out.Bounds().Dy())
}

func (s *BildProcessorSuite) TestBildProcessor_Grayscale() {
	var actual, expected []byte
	var err error
	out := s.processor.GrayScale(s.srcImage)
	actual, err = s.processor.Encode(out, "png")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_grayscaled.png")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)

	assert.EqualValues(s.T(), actual, expected)
}

func (s *BildProcessorSuite) TestBildProcessor_Flip() {
	var actual, expected []byte
	var err error
	out := s.processor.Flip(s.srcImage, "v")
	actual, err = s.processor.Encode(out, "jpeg")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_flipedV.jpg")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), actual, expected)

	out = s.processor.Flip(s.srcImage, "h")
	actual, err = s.processor.Encode(out, "jpeg")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_flipedH.jpg")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), actual, expected)

	out = s.processor.Flip(s.srcImage, "vh")
	actual, err = s.processor.Encode(out, "jpeg")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_flipedVH.jpg")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), actual, expected)
}

func (s *BildProcessorSuite) TestBildProcessor_Rotate() {
	var actual, expected []byte
	var err error
	out := s.processor.Rotate(s.srcImage, 90)
	actual, err = s.processor.Encode(out, "jpeg")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_rotated90.jpg")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), actual, expected)

	out = s.processor.Rotate(s.srcImage, 175)
	actual, err = s.processor.Encode(out, "jpeg")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_rotated175.jpg")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), actual, expected)

	out = s.processor.Rotate(s.srcImage, 450)
	actual, err = s.processor.Encode(out, "jpeg")
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_rotated90.jpg")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)
	assert.EqualValues(s.T(), actual, expected)
}

func (s *BildProcessorSuite) TestBildProcessor_Watermark() {
	output, err := s.processor.Watermark(s.srcData, s.watermarkData, 200)

	assert.NotNil(s.T(), output)
	assert.Nil(s.T(), err)

	assert.NotEqual(s.T(), s.srcData, output)
}

func (s *BildProcessorSuite) TestBildProcessorWithBadInput() {
	output, err := s.processor.Watermark(s.badData, s.watermarkData, 255)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)

	output, err = s.processor.Watermark(s.srcData, s.badData, 255)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)
}
