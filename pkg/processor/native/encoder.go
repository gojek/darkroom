package native

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
)

var DefaultEncoderOptions = &EncoderOptions{
	JpegQuality:         jpeg.DefaultQuality,
	PngCompressionLevel: png.BestCompression,
}

type EncoderOptions struct {
	JpegQuality         int
	PngCompressionLevel png.CompressionLevel
}

type Encoder interface {
	Encode(img image.Image) ([]byte, error)
}

type JpegEncoder struct {
	Option *jpeg.Options
}
type PngEncoder struct {
	Encoder *png.Encoder
}
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

type Encoders struct {
	options     *EncoderOptions
	jpegEncoder Encoder
	pngEncoder  Encoder
	noOpEncoder Encoder
}

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

func (e *Encoders) JpegEncoder() Encoder {
	return e.jpegEncoder
}

func (e *Encoders) PngEncoder() Encoder {
	return e.pngEncoder
}

func (e *Encoders) Options() *EncoderOptions {
	return e.options
}

func NewEncoders(opts *EncoderOptions) *Encoders {
	return &Encoders{
		options:     opts,
		jpegEncoder: &JpegEncoder{Option: &jpeg.Options{Quality: opts.JpegQuality}},
		pngEncoder: &PngEncoder{
			Encoder: &png.Encoder{CompressionLevel: opts.PngCompressionLevel},
		},
		noOpEncoder: &NopEncoder{},
	}
}
