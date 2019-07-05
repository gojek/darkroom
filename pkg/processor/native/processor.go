package native

import (
	"bytes"
	"github.com/anthonynsimon/bild/clone"
	"github.com/anthonynsimon/bild/parallel"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gojek/darkroom/pkg/metrics"
	"github.com/gojek/darkroom/pkg/processor"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"time"
)

const (
	pngType = "png"
	jpgType = "jpeg"

	watermarkDurationKey = "watermarkDuration"
	decodeDurationKey    = "decodeDuration"
	encodeDurationKey    = "encodeDuration"
)

// BildProcessor uses bild library to process images using native Golang image.Image interface
type BildProcessor struct {
}

// Crop takes an input byte array, width, height and a CropPoint and returns the cropped image bytes or error
func (bp *BildProcessor) Crop(img image.Image, width, height int, point processor.CropPoint) image.Image {
	w, h := getResizeWidthAndHeightForCrop(width, height, img.Bounds().Dx(), img.Bounds().Dy())

	img = transform.Resize(img, w, h, transform.Linear)
	x0, y0 := getStartingPointForCrop(w, h, width, height, point)
	rect := image.Rect(x0, y0, width+x0, height+y0)
	img = (clone.AsRGBA(img)).SubImage(rect)

	return img
}

// Resize takes an input byte array, width and height and returns the re-sized image bytes or error
func (bp *BildProcessor) Resize(img image.Image, width, height int) image.Image {

	initW := img.Bounds().Dx()
	initH := img.Bounds().Dy()

	w, h := getResizeWidthAndHeight(width, height, initW, initH)
	if w != initW || h != initH {
		img = transform.Resize(img, w, h, transform.Linear)
	}

	return img
}

// Watermark takes an input byte array, overlay byte array and opacity value
// and returns the watermarked image bytes or error
func (bp *BildProcessor) Watermark(base []byte, overlay []byte, opacity uint8) ([]byte, error) {
	baseImg, f, err := bp.Decode(base)
	if err != nil {
		return nil, err
	}
	overlayImg, _, err := bp.Decode(overlay)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	ratio := float64(overlayImg.Bounds().Dy()) / float64(overlayImg.Bounds().Dx())
	dWidth := float64(baseImg.Bounds().Dx()) / 2.0

	// Resizing overlay image according to base image
	overlayImg = transform.Resize(overlayImg, int(dWidth), int(dWidth*ratio), transform.Linear)

	// Anchor point for overlaying
	x := (baseImg.Bounds().Dx() - overlayImg.Bounds().Dx()) / 2
	y := (baseImg.Bounds().Dy() - overlayImg.Bounds().Dy()) / 2
	offset := image.Pt(int(x), int(y))

	// Mask image (that is just a solid light gray image)
	mask := image.NewUniform(color.Alpha{A: opacity})

	// Performing overlay
	draw.DrawMask(baseImg.(draw.Image), overlayImg.Bounds().Add(offset), overlayImg, image.ZP, mask, image.ZP, draw.Over)
	metrics.Update(metrics.UpdateOption{Name: watermarkDurationKey, Type: metrics.Duration, Duration: time.Since(t)})

	return bp.Encode(baseImg, f)
}

// GrayScale takes an input byte array and returns the grayscaled byte array or error
func (bp *BildProcessor) GrayScale(img image.Image) image.Image {
	src := clone.AsRGBA(img)
	bounds := src.Bounds()
	if bounds.Empty() {
		src = &image.RGBA{}
	} else {
		parallel.Line(bounds.Dy(), func(start, end int) {
			for y := start; y < end; y++ {
				for x := 0; x < bounds.Dx(); x++ {
					srcPix := src.At(x, y).(color.RGBA)
					g := color.GrayModel.Convert(srcPix).(color.Gray).Y
					src.Set(x, y, color.RGBA{R: g, G: g, B: g, A: srcPix.A})
				}
			}
		})
	}
	return src
}

func (bp *BildProcessor) Decode(data []byte) (image.Image, string, error) {
	t := time.Now()
	img, f, err := image.Decode(bytes.NewReader(data))
	if err == nil {
		metrics.Update(metrics.UpdateOption{Name: decodeDurationKey, Type: metrics.Duration, Duration: time.Since(t)})
	}
	return img, f, err
}

func (bp *BildProcessor) Encode(img image.Image, format string) ([]byte, error) {
	t := time.Now()
	if format == pngType && isOpaque(img) {
		format = jpgType
	}
	buff := &bytes.Buffer{}
	var err error
	if format == pngType {
		enc := png.Encoder{CompressionLevel: png.BestCompression}
		err = enc.Encode(buff, img)
	} else {
		err = jpeg.Encode(buff, img, nil)
	}
	metrics.Update(metrics.UpdateOption{Name: encodeDurationKey, Type: metrics.Duration, Duration: time.Since(t)})
	return buff.Bytes(), err
}

// NewBildProcessor creates a new BildProcessor
func NewBildProcessor() *BildProcessor {
	return &BildProcessor{}
}
