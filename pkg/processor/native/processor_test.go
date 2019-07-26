package native

import (
	"image"
	"image/png"
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

func (s *BildProcessorSuite) TestNewBildProcessorWithCompression() {
	p := NewBildProcessorWithCompression(&CompressionOptions{JpegQuality: 70, PngCompressionLevel: png.BestSpeed})

	assert.NotNil(s.T(), p)
	assert.Equal(s.T(), 70, p.encoders.Options().JpegQuality)
	assert.Equal(s.T(), png.BestSpeed, p.encoders.Options().PngCompressionLevel)
}

func (s *BildProcessorSuite) TestBildProcessor_Resize() {
	out := s.processor.Resize(s.srcImage, 600, 500)

	assert.NotNil(s.T(), out)
	assert.Equal(s.T(), 600, out.Bounds().Dx())
	assert.Equal(s.T(), 450, out.Bounds().Dy())
}

func (s *BildProcessorSuite) TestBildProcessor_ResizeWithSameWidthAndHeight() {
	out := s.processor.Resize(s.srcImage, 500, 375)

	assert.NotNil(s.T(), out)
	assert.Equal(s.T(), 500, out.Bounds().Dx())
	assert.Equal(s.T(), 375, out.Bounds().Dy())
	// Checks if the image is the same image which was passed by doing a pointer comparision
	assert.Equal(s.T(), &s.srcImage, &out)
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
	cases := []struct {
		flipMode string
		testFile string
	}{
		{
			flipMode: "v",
			testFile: "_testdata/test_flipedV.jpg",
		},
		{
			flipMode: "h",
			testFile: "_testdata/test_flipedH.jpg",
		},
		{
			flipMode: "vh",
			testFile: "_testdata/test_flipedVH.jpg",
		},
	}

	for _, c := range cases {
		out := s.processor.Flip(s.srcImage, c.flipMode)
		actual, err = s.processor.Encode(out, "jpeg")
		assert.NotNil(s.T(), actual)
		assert.Nil(s.T(), err)
		expected, err = ioutil.ReadFile(c.testFile)
		assert.NotNil(s.T(), expected)
		assert.Nil(s.T(), err)
		assert.EqualValues(s.T(), actual, expected)
	}
}

func (s *BildProcessorSuite) TestBildProcessor_Rotate() {
	var actual, expected []byte
	var err error
	cases := []struct {
		angle    float64
		testFile string
	}{
		{
			angle:    90.0,
			testFile: "_testdata/test_rotated90.jpg",
		},
		{
			angle:    175,
			testFile: "_testdata/test_rotated175.jpg",
		},
		{
			angle:    450.0,
			testFile: "_testdata/test_rotated90.jpg",
		},
	}

	for _, c := range cases {
		out := s.processor.Rotate(s.srcImage, c.angle)
		actual, err = s.processor.Encode(out, "jpeg")
		assert.NotNil(s.T(), actual)
		assert.Nil(s.T(), err)
		expected, err = ioutil.ReadFile(c.testFile)
		assert.NotNil(s.T(), expected)
		assert.Nil(s.T(), err)
		assert.EqualValues(s.T(), actual, expected)
	}
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
