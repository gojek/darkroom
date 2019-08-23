package native

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
	"github.com/gojek/darkroom/pkg/processor"
)

// Encoder is an interface to Encode image and return the encoded byte array or error
type Encoder interface {
	Encode(img image.Image) ([]byte, error)
}

// JPEGEncoder is an object to encode image to byte array with jpeg format
type JPEGEncoder struct {
	Options *jpeg.Options
}

// PNGEncoder is an object to encode image to byte array with png format
type PNGEncoder struct {
	Encoder *png.Encoder
}

// WebPEncoder is an object to encode image to byte array with webp format
type WebPEncoder struct {
	Options *webp.Options
}

// NoOpEncoder is a no-op encoder object for unsupported format and will return error
type NoOpEncoder struct{}

func (e *PNGEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := e.Encoder.Encode(buff, img)
	return buff.Bytes(), err
}

func (e *JPEGEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := jpeg.Encode(buff, img, e.Options)
	return buff.Bytes(), err
}

func (e *WebPEncoder) Encode(img image.Image) ([]byte, error) {
	buff := &bytes.Buffer{}
	err := webp.Encode(buff, img, e.Options)
	return buff.Bytes(), err
}

func (e *NoOpEncoder) Encode(img image.Image) ([]byte, error) {
	return nil, errors.New("unknown format: failed to encode image")
}

// Encoders is a struct to store all supported encoders so that we don't have to create new encoder every time
type Encoders struct {
	JPEGEncoder *JPEGEncoder
	PNGEncoder  *PNGEncoder
	NoOpEncoder *NoOpEncoder
	WebPEncoder *WebPEncoder
}

// EncodersOption represents builder function for Encoders
type EncodersOption func(*Encoders)

// GetEncoder takes an input of image and extension and return the appropriate Encoder for encoding the image
func (e *Encoders) GetEncoder(img image.Image, format string) Encoder {
	switch format {
	case processor.FormatJPG, processor.FormatJPEG:
		return e.JPEGEncoder
	case processor.FormatPNG:
		if e.JPEGEncoder.Options.Quality != 100 && isOpaque(img) {
			return e.JPEGEncoder
		}
		return e.PNGEncoder
	case processor.FormatWebP:
		return e.WebPEncoder
	default:
		return e.NoOpEncoder
	}
}

// Support takes an input of image format and return whether encoders support encoding for that image format
func (e *Encoders) Support(format string) bool {
	switch format {
	case processor.FormatJPG, processor.FormatJPEG, processor.FormatPNG, processor.FormatWebP:
		return true
	default:
		return false
	}
}

// WithJPEGEncoder is a builder function for setting custom JPEGEncoder
func WithJPEGEncoder(j *JPEGEncoder) EncodersOption {
	return func(e *Encoders) {
		e.JPEGEncoder = j
	}
}

// WithPNGEncoder is a builder function for setting custom PNGEncoder
func WithPNGEncoder(p *PNGEncoder) EncodersOption {
	return func(e *Encoders) {
		e.PNGEncoder = p
	}
}

// WithWebPEncoder is a builder function for setting custom WebPEncoder
func WithWebPEncoder(w *WebPEncoder) EncodersOption {
	return func(e *Encoders) {
		e.WebPEncoder = w
	}
}

// NewEncoders creates a new Encoders, if called without parameter (builder), all encoders option will be default
func NewEncoders(opts ...EncodersOption) *Encoders {
	e := &Encoders{
		JPEGEncoder: &JPEGEncoder{Options: &jpeg.Options{Quality: jpeg.DefaultQuality}},
		PNGEncoder: &PNGEncoder{
			Encoder: &png.Encoder{CompressionLevel: png.BestCompression},
		},
		NoOpEncoder: &NoOpEncoder{},
		WebPEncoder: &WebPEncoder{},
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}
