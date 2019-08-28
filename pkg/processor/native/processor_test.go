package native

import (
	"bytes"
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

func TestBildProcessor(t *testing.T) {
	suite.Run(t, new(BildProcessorSuite))
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
	actual, err = s.processor.Encode(out, "png", false)
	assert.NotNil(s.T(), actual)
	assert.Nil(s.T(), err)
	expected, err = ioutil.ReadFile("_testdata/test_grayscaled.png")
	assert.NotNil(s.T(), expected)
	assert.Nil(s.T(), err)

	assert.EqualValues(s.T(), actual, expected)
}

func (s *BildProcessorSuite) TestBildProcessor_Blur() {
	var actual, expected []byte
	var err error
	cases := []struct {
		radius       float64
		expectedFile string
	}{
		{
			radius:       0.0,
			expectedFile: "_testdata/test.jpg",
		},
		{
			radius:       1.0,
			expectedFile: "_testdata/test_blurred_1.jpg",
		},
		{
			radius:       60.0,
			expectedFile: "_testdata/test_blurred_60.jpg",
		},
	}
	for _, c := range cases {
		out := s.processor.Blur(s.srcImage, c.radius)
		actual, err = s.processor.Encode(out, "jpeg", false)
		assert.NotNil(s.T(), actual)
		assert.Nil(s.T(), err)
		expected, err = ioutil.ReadFile(c.expectedFile)
		assert.NotNil(s.T(), expected)
		assert.Nil(s.T(), err)
		assert.EqualValues(s.T(), actual, expected)
	}
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
		actual, err = s.processor.Encode(out, "jpeg", false)
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
		actual, err = s.processor.Encode(out, "jpeg", false)
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

func (s *BildProcessorSuite) TestBildProcessor_FixOrientation() {
	var testFiles = []string{
		"./_testdata/exif_orientation/f2t.jpg",
		"./_testdata/exif_orientation/f3t.jpg",
		"./_testdata/exif_orientation/f4t.jpg",
		"./_testdata/exif_orientation/f5t.jpg",
		"./_testdata/exif_orientation/f6t.jpg",
		"./_testdata/exif_orientation/f7t.jpg",
		"./_testdata/exif_orientation/f8t.jpg",
	}
	expected, err := ioutil.ReadFile("./_testdata/exif_orientation/expected.jpg")
	if err != nil {
		panic(err)
	}
	for _, testFile := range testFiles {
		file, err := ioutil.ReadFile(testFile)
		if err != nil {
			panic(err)
		}
		orientation, _ := GetOrientation(bytes.NewReader(file))
		img, _, err := s.processor.Decode(file)
		img = s.processor.FixOrientation(img, orientation)
		actual, err := s.processor.Encode(img, "jpg", false)
		assert.Nil(s.T(), err)
		assert.EqualValues(s.T(), expected, actual)
	}
}

func (s *BildProcessorSuite) TestBildProcessor_Watermark_WithBadInput() {
	output, err := s.processor.Watermark(s.badData, s.watermarkData, 255)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)

	output, err = s.processor.Watermark(s.srcData, s.badData, 255)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), output)
}

func (s *BildProcessorSuite) TestBildProcessor_WithEncoders() {
	e := NewEncoders()
	bp := NewBildProcessor(WithEncoders(e))
	assert.Equal(s.T(), e, bp.encoders)
}

func (s *BildProcessorSuite) TestBildProcessor_Decode_GivenWebPImageShouldBeAbleToDecodeProperly() {
	data, _ := ioutil.ReadFile("_testdata/test.webp")
	_, ext, err := s.processor.Decode(data)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "webp", ext)
}

func (s *BildProcessorSuite) TestBildProcessor_Support_GivenSupportedFormatShouldReturnTrue() {
	assert.True(s.T(), s.processor.Support(processor.FormatJPG))
	assert.True(s.T(), s.processor.Support(processor.FormatJPEG))
	assert.True(s.T(), s.processor.Support(processor.FormatPNG))
	assert.True(s.T(), s.processor.Support(processor.FormatWebP))
}

func (s *BildProcessorSuite) TestBildProcessor_Support_GivenUnsupportedFormatShouldReturnFalse() {
	assert.False(s.T(), s.processor.Support("unknown"))
}
