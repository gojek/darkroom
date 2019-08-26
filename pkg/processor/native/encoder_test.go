package native

import (
	"bytes"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"testing"

	"github.com/gojek/darkroom/pkg/processor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type EncoderSuite struct {
	suite.Suite
	encoders         *Encoders
	srcImage         image.Image
	processor        processor.Processor
	opaqueImage      image.Image
	transparentImage image.Image
}

func (s *EncoderSuite) SetupSuite() {
	s.encoders = NewEncoders()

	s.processor = NewBildProcessor()
	data, err := ioutil.ReadFile("_testdata/test.png")
	if err != nil {
		panic(err)
	}

	s.srcImage, _, err = image.Decode(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}

	opaque := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(opaque, opaque.Bounds(), image.Opaque, image.ZP, draw.Src)
	s.opaqueImage = opaque

	transparent := image.NewRGBA(image.Rect(0, 0, 640, 480))
	draw.Draw(transparent, transparent.Bounds(), image.Transparent, image.ZP, draw.Src)
	s.transparentImage = transparent
}

func TestEncoder(t *testing.T) {
	suite.Run(t, new(EncoderSuite))
}

func TestNewEncoders(t *testing.T) {
	jpegEncoder := &JPEGEncoder{}
	pngEncoder := &PNGEncoder{}
	webPEncoder := &WebPEncoder{}
	e := NewEncoders(
		WithJPEGEncoder(jpegEncoder),
		WithPNGEncoder(pngEncoder),
		WithWebPEncoder(webPEncoder),
	)
	assert.Equal(t, jpegEncoder, e.JPEGEncoder)
	assert.Equal(t, pngEncoder, e.PNGEncoder)
	assert.Equal(t, webPEncoder, e.WebPEncoder)
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenJpgExtensionShouldReturnJpegEncoder() {
	assert.IsType(s.T(), &JPEGEncoder{}, s.encoders.GetEncoder(s.opaqueImage, processor.FormatJPG, false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenJpegExtensionShouldReturnJpegEncoder() {
	assert.IsType(s.T(), &JPEGEncoder{}, s.encoders.GetEncoder(s.opaqueImage, processor.FormatJPEG, false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenOpaqueImageAndPngExtensionShouldReturnPngEncoder() {
	s.encoders.JPEGEncoder.Options.Quality = 99
	assert.IsType(s.T(), &JPEGEncoder{}, s.encoders.GetEncoder(s.opaqueImage, processor.FormatPNG, false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenOpaqueImageAndPngExtensionShouldReturnJpegEncoder() {
	s.encoders.JPEGEncoder.Options.Quality = 100
	assert.IsType(s.T(), &PNGEncoder{}, s.encoders.GetEncoder(s.opaqueImage, processor.FormatPNG, false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenTransparentImageAndPngExtensionShouldReturnPngEncoder() {
	assert.IsType(s.T(), &PNGEncoder{}, s.encoders.GetEncoder(s.transparentImage, processor.FormatPNG, false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenUnknownExtensionShouldReturnNopEncoder() {
	assert.IsType(s.T(), &NoOpEncoder{}, s.encoders.GetEncoder(image.Black, "unknown", false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenWebPExtensionShouldReturnWebPEncoder() {
	assert.IsType(s.T(), &WebPEncoder{}, s.encoders.GetEncoder(s.transparentImage, processor.FormatWebP, false))
}

func (s *EncoderSuite) TestEncoders_GetEncoder_GivenEnforceFmtTrueShouldReturnCorrectEncoder() {
	s.encoders.JPEGEncoder.Options.Quality = 100
	assert.IsType(s.T(), &PNGEncoder{}, s.encoders.GetEncoder(s.opaqueImage, processor.FormatPNG, true))
}

func (s *EncoderSuite) TestJpgEncoder_Encode_ShouldEncodeToJpeg() {
	encoder := JPEGEncoder{Options: nil}
	data, err := encoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)
	_, f, err := NewBildProcessor().Decode(data)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), "jpeg", f)
}

func (s *EncoderSuite) TestJpgEncoder_Encode_QualityShouldAffectFileSize() {
	lowQualityEncoder := JPEGEncoder{Options: &jpeg.Options{Quality: 25}}
	highQualityEncoder := JPEGEncoder{Options: &jpeg.Options{Quality: 90}}
	lowQualityData, err := lowQualityEncoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)

	highQualityData, err := highQualityEncoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)

	assert.True(s.T(), len(lowQualityData) < len(highQualityData))
}

func (s *EncoderSuite) TestNopEncoder() {
	nopEncoder := NoOpEncoder{}

	data, err := nopEncoder.Encode(s.srcImage)
	assert.Nil(s.T(), data)
	assert.Error(s.T(), err)
}

func (s *EncoderSuite) TestPngEncoder_Encode_ShouldEncodeToPng() {
	encoder := PNGEncoder{Encoder: &png.Encoder{CompressionLevel: png.NoCompression}}
	data, err := encoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)
	_, f, err := s.processor.Decode(data)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), processor.FormatPNG, f)
}

func (s *EncoderSuite) TestWebPEncoder_Encode_ShouldEncodeToWebP() {
	encoder := WebPEncoder{}
	data, err := encoder.Encode(s.srcImage)
	assert.Nil(s.T(), err)
	_, f, err := s.processor.Decode(data)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), processor.FormatWebP, f)
}
