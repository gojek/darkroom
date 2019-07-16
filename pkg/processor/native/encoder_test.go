package native

import (
	"bytes"
	"github.com/gojek/darkroom/pkg/processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"testing"
)

type EncoderSuite struct {
	suite.Suite
	srcImage  image.Image
	processor processor.Processor
}

func (s *EncoderSuite) SetupSuite() {
	s.processor = NewBildProcessor()
	data, err := ioutil.ReadFile("_testdata/test.png")
	if err != nil {
		panic(err)
	}
	s.srcImage, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
}

func TestNewEncoders(t *testing.T) {
	suite.Run(t, new(EncoderSuite))
}

func (s *EncoderSuite) TestEncoders_GetEncoder() {
	encoders := NewEncoders(DefaultCompressionOptions)
	_, ok := (encoders.GetEncoder(image.Opaque, "jpg")).(*JpegEncoder)
	assert.True(s.T(), ok)
	_, ok = (encoders.GetEncoder(image.Opaque, "jpeg")).(*JpegEncoder)
	assert.True(s.T(), ok)
	encoders.options.JpegQuality = 99
	_, ok = (encoders.GetEncoder(image.Opaque, "png")).(*JpegEncoder)
	assert.True(s.T(), ok)
	encoders.options.JpegQuality = 100
	_, ok = (encoders.GetEncoder(image.Opaque, "png")).(*PngEncoder)
	assert.True(s.T(), ok)
	_, ok = (encoders.GetEncoder(image.Transparent, "png")).(*PngEncoder)
	assert.True(s.T(), ok)
	_, ok = (encoders.GetEncoder(image.Black, "unknown")).(*NopEncoder)
	assert.True(s.T(), ok)
}

func (s *EncoderSuite) TestJpgEncoder_Encode_ShouldEncodeToJpeg() {
	encoder := JpegEncoder{Option: nil}
	data, err := encoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)
	_, f, err := NewBildProcessor().Decode(data)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "jpeg", f)
}

func (s *EncoderSuite) TestJpgEncoder_Encode_QualityShouldAffectFileSize() {
	lowQualityEncoder := JpegEncoder{Option: &jpeg.Options{Quality: 25}}
	highQualityEncoder := JpegEncoder{Option: &jpeg.Options{Quality: 90}}
	lowQualityData, err := lowQualityEncoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)

	highQualityData, err := highQualityEncoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)

	assert.True(s.T(), len(lowQualityData) < len(highQualityData))
}

func (s *EncoderSuite) TestPngEncoder_Encode_ShouldEncodeToPng() {
	encoder := PngEncoder{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}
	data, err := encoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)
	_, f, err := s.processor.Decode(data)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "png", f)
}
