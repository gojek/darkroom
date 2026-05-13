package service

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gojek/darkroom/pkg/metrics"
	"github.com/gojek/darkroom/pkg/processor"
	"github.com/gojek/darkroom/pkg/processor/native"
)

const (
	width        = "w"
	height       = "h"
	fit          = "fit"
	crop         = "crop"
	mono         = "mono"
	blackHexCode = "000000"
	flip         = "flip"
	rotate       = "rot"
	auto         = "auto"
	blur         = "blur"
	compress     = "compress"
	format       = "format"
	scale        = "scale"
	quantize     = "quantize"

	cropDurationKey      = "cropDuration"
	decodeDurationKey    = "decodeDuration"
	encodeDurationKey    = "encodeDuration"
	grayScaleDurationKey = "grayScaleDuration"
	blurDurationKey      = "blurDuration"
	resizeDurationKey    = "resizeDuration"
	flipDurationKey      = "flipDuration"
	rotateDurationKey    = "rotateDuration"
	fixOrientationKey    = "fixOrientation"
	scaleDurationKey     = "scaleDuration"
)

var plan9WithTransparency color.Palette

func init() {
	cp := make(color.Palette, len(palette.Plan9))
	copy(cp, palette.Plan9)
	cp[len(cp)-1] = color.NRGBA{0, 0, 0, 0}
	plan9WithTransparency = cp
}

// Manipulator interface sets the contract on the implementation for common processing support in darkroom
type Manipulator interface {
	// Process takes ProcessSpec as an argument and returns []byte, error
	Process(spec processSpec) ([]byte, error)

	// HasDefaultParams returns true if defaultParams are present, returns false otherwise
	HasDefaultParams() bool
}

type manipulator struct {
	processor     processor.Processor
	defaultParams map[string]string
	metricService metrics.MetricService
}

// Process takes ProcessSpec as an argument and returns []byte, error
// This manipulator uses bild to do the actual image manipulations
func (m *manipulator) Process(spec processSpec) ([]byte, error) {
	params := spec.Params
	params = joinParams(params, m.defaultParams)
	var err error
	t := time.Now()
	data, f, err := m.processor.Decode(spec.ImageData)
	if err != nil {
		return nil, err
	}

	var originalPalette color.Palette
	if paletted, ok := data.(*image.Paletted); ok {
		originalPalette = paletted.Palette
	}

	if spec.TargetFormat != "" {
		f = spec.TargetFormat
	}

	m.metricService.TrackDuration(decodeDurationKey, t, spec.ImageData)
	if params[fit] == crop {
		t = time.Now()
		data = m.processor.Crop(data, CleanInt(params[width]), CleanInt(params[height]), GetCropPoint(params[crop]))
		m.metricService.TrackDuration(cropDurationKey, t, spec.ImageData)
	} else if params[fit] == scale {
		t = time.Now()
		data = m.processor.Scale(data, CleanInt(params[width]), CleanInt(params[height]))
		m.metricService.TrackDuration(scaleDurationKey, t, spec.ImageData)
	} else if len(params[fit]) == 0 && (CleanInt(params[width]) != 0 || CleanInt(params[height]) != 0) {
		t = time.Now()
		data = m.processor.Resize(data, CleanInt(params[width]), CleanInt(params[height]))
		m.metricService.TrackDuration(resizeDurationKey, t, spec.ImageData)
	}

	if params[mono] == blackHexCode {
		t = time.Now()
		data = m.processor.GrayScale(data)
		m.metricService.TrackDuration(grayScaleDurationKey, t, spec.ImageData)
	}
	if radius := CleanFloat(params[blur], 1000); radius > 0 {
		t = time.Now()
		data = m.processor.Blur(data, radius)
		m.metricService.TrackDuration(blurDurationKey, t, spec.ImageData)
	}

	if len(params[flip]) != 0 {
		t = time.Now()
		data = m.processor.Flip(data, params[flip])
		m.metricService.TrackDuration(flipDurationKey, t, spec.ImageData)
	}

	if angle := CleanFloat(params[rotate], 360); angle > 0 {
		t = time.Now()
		data = m.processor.Rotate(data, angle)
		m.metricService.TrackDuration(rotateDurationKey, t, spec.ImageData)
	}

	if params[quantize] == "true" && strings.EqualFold(f, processor.ExtensionPNG) {
		if _, isPal := data.(*image.Paletted); !isPal {
			quantizePalette := palette.Plan9
			if originalPalette != nil {
				quantizePalette = originalPalette
			} else if imageHasTransparency(data) {
				quantizePalette = plan9WithTransparency
			}
			data = convertToPaletted(data, quantizePalette)
		}
	}

	autos := strings.Split(params[auto], ",")
	originalFormat := f
	for _, a := range autos {
		if a == compress {
			orientation, _ := native.GetOrientation(bytes.NewReader(spec.ImageData))
			t = time.Now()
			data = m.processor.FixOrientation(data, orientation)
			m.metricService.TrackDuration(fixOrientationKey, t, spec.ImageData)
		} else if a == format {
			w := spec.IsWebPSupported()
			if w {
				f = processor.ExtensionWebP
			} else if f == processor.ExtensionWebP {
				f = processor.ExtensionPNG
			}
		}
	}

	t = time.Now()
	src, err := m.processor.Encode(data, f)
	if err != nil && f != originalFormat {
		src, err = m.processor.Encode(data, originalFormat)
	}
	if err == nil {
		m.metricService.TrackDuration(encodeDurationKey, t, spec.ImageData)
	}
	return src, err
}

func convertToPaletted(img image.Image, p color.Palette) image.Image {
	if _, ok := img.(*image.Paletted); ok {
		return img
	}
	bounds := img.Bounds()
	palettedImg := image.NewPaletted(bounds, p)
	draw.FloydSteinberg.Draw(palettedImg, bounds, img, bounds.Min)

	transparentIndex := paletteTransparentIndex(palettedImg.Palette)
	if transparentIndex >= 0 {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				_, _, _, a := img.At(x, y).RGBA()
				if a == 0 {
					palettedImg.SetColorIndex(x, y, uint8(transparentIndex))
				}
			}
		}
	}
	return palettedImg
}

func imageHasTransparency(img image.Image) bool {
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a < 0xffff {
				return true
			}
		}
	}
	return false
}

func paletteTransparentIndex(p color.Palette) int {
	for i, c := range p {
		_, _, _, a := c.RGBA()
		if a == 0 {
			return i
		}
	}
	return -1
}

func paletteHasTransparency(p color.Palette) bool {
	return paletteTransparentIndex(p) >= 0
}

// HasDefaultParams returns true if defaultParams are present, returns false otherwise
func (m *manipulator) HasDefaultParams() bool {
	return len(m.defaultParams) > 0
}

func joinParams(params map[string]string, defaultParams map[string]string) map[string]string {
	fp := make(map[string]string)
	for p := range defaultParams {
		fp[p] = defaultParams[p]
	}
	for p := range params {
		if fp[p] != "" {
			fp[p] = fmt.Sprintf("%s,%s", defaultParams[p], params[p])
		} else {
			fp[p] = params[p]
		}
	}
	return fp
}

// CleanInt takes a string and return an int not greater than 9999
func CleanInt(input string) int {
	val, _ := strconv.Atoi(input)
	if val <= 0 {
		return 0
	}
	return val % 10000 // Never return value greater than 9999
}

// CleanFloat takes a string and return a float64 not greater than bound
func CleanFloat(input string, bound float64) float64 {
	val, _ := strconv.ParseFloat(input, 64)
	if val <= 0 {
		return 0
	}
	return math.Mod(val, bound) // Never return value greater than bound
}

// GetCropPoint takes a string and returns the type Point
func GetCropPoint(input string) processor.Point {
	switch input {
	case "top":
		return processor.PointTop
	case "top,left":
		return processor.PointTopLeft
	case "top,right":
		return processor.PointTopRight
	case "left":
		return processor.PointLeft
	case "right":
		return processor.PointRight
	case "bottom":
		return processor.PointBottom
	case "bottom,left":
		return processor.PointBottomLeft
	case "bottom,right":
		return processor.PointBottomRight
	default:
		return processor.PointCenter
	}
}

// NewManipulator takes in a Processor interface and returns a new Manipulator
func NewManipulator(processor processor.Processor, defaultParams map[string]string,
	metricService metrics.MetricService) Manipulator {
	return &manipulator{
		processor:     processor,
		defaultParams: defaultParams,
		metricService: metricService,
	}
}
