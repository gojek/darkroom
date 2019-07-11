package native

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
)

// DefaultCompressionOptions is the default compression options used for encoding images
var DefaultCompressionOptions = &CompressionOptions{
	JpegQuality:         jpeg.DefaultQuality,
	PngCompressionLevel: png.BestCompression,
}

// CompressionOptions is an object to configure jpeg quality and png compression level when encoding the image
type CompressionOptions struct {
	JpegQuality         int
	PngCompressionLevel png.CompressionLevel
}

// Encoder is an interface to Encode image and return the encoded byte array or error
type Encoder interface {
	Encode(img image.Image) ([]byte, error)
}

// JpegEncoder is an object to encode image to byte array with jpeg format
type JpegEncoder struct {
	Option *jpeg.Options
}

// PngEncoder is an object to encode image to byte array with png format
type PngEncoder struct {
	Encoder *png.Encoder
}

// NopEncoder is a no-op encoder object for unsupported format and will return error
type NopEncoder struct{}

func (e *PngEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := e.Encoder.Encode(buff, img)
	return buff.Bytes(), err
}

func (e *JpegEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := jpeg.Encode(buff, img, e.Option)
	return buff.Bytes(), err
}

func (e *NopEncoder) Encode(img image.Image) ([]byte, error) {
	return nil, errors.New("unknown format: failed to encode image")
}

// Encoders is a struct to store all supported encoders so that we don't have to create new encoder every time
type Encoders struct {
	options     *CompressionOptions
	jpegEncoder Encoder
	pngEncoder  Encoder
	noOpEncoder Encoder
}

// GetEncoder takes an input of image and extension and return the appropriate Encoder for encoding the image
func (e *Encoders) GetEncoder(img image.Image, ext string) Encoder {
	if ext == "jpg" || ext == "jpeg" {
		return e.jpegEncoder
	}
	if ext == "png" {
		if isOpaque(img) {
			return e.jpegEncoder
		}
		return e.pngEncoder
	}
	return e.noOpEncoder
}

// Getter for JpegEncoder
func (e *Encoders) JpegEncoder() Encoder {
	return e.jpegEncoder
}

// Getter for PngEncoder
func (e *Encoders) PngEncoder() Encoder {
	return e.pngEncoder
}

// Getter for Options
func (e *Encoders) Options() *CompressionOptions {
	return e.options
}

func NewEncoders(opts *CompressionOptions) *Encoders {
	return &Encoders{
		options:     opts,
		jpegEncoder: &JpegEncoder{Option: &jpeg.Options{Quality: opts.JpegQuality}},
		pngEncoder: &PngEncoder{
			Encoder: &png.Encoder{CompressionLevel: opts.PngCompressionLevel},
		},
		noOpEncoder: &NopEncoder{},
	}
}
