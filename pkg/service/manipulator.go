package service

import (
	"bytes"
	"fmt"
	"math"
	"net/http"
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

// Manipulator interface sets the contract on the implementation for common processing support in darkroom
type Manipulator interface {
	// Process takes ProcessSpec as an argument and returns []byte, error
	Process(spec processSpec) ([]byte, error)
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
	m.metricService.TrackDecodeDuration(t, spec.ImageData)
	if params[fit] == crop {
		t = time.Now()
		data = m.processor.Crop(data, CleanInt(params[width]), CleanInt(params[height]), GetCropPoint(params[crop]))
		m.metricService.TrackCropDuration(t, spec.ImageData)
	} else if params[fit] == scale {
		t = time.Now()
		data = m.processor.Scale(data, CleanInt(params[width]), CleanInt(params[height]))
		m.metricService.TrackScaleDuration(t, spec.ImageData)
	} else if len(params[fit]) == 0 && (CleanInt(params[width]) != 0 || CleanInt(params[height]) != 0) {
		t = time.Now()
		data = m.processor.Resize(data, CleanInt(params[width]), CleanInt(params[height]))
		m.metricService.TrackResizeDuration(t, spec.ImageData)
	}

	if params[mono] == blackHexCode {
		t = time.Now()
		data = m.processor.GrayScale(data)
		m.metricService.TrackGrayScaleDuration(t, spec.ImageData)
	}
	if radius := CleanFloat(params[blur], 1000); radius > 0 {
		t = time.Now()
		data = m.processor.Blur(data, radius)
		m.metricService.TrackBlurDuration(t, spec.ImageData)
	}

	autos := strings.Split(params[auto], ",")
	for _, a := range autos {
		if a == compress {
			orientation, _ := native.GetOrientation(bytes.NewReader(spec.ImageData))
			t = time.Now()
			data = m.processor.FixOrientation(data, orientation)
			m.metricService.TrackFixOrientationDuration(t, spec.ImageData)
		} else if a == format {
			w := spec.IsWebPSupported()
			if w {
				f = processor.ExtensionWebP
			} else if f == processor.ExtensionWebP {
				f = processor.ExtensionPNG
			}
		}
	}

	if len(params[flip]) != 0 {
		t = time.Now()
		data = m.processor.Flip(data, params[flip])
		m.metricService.TrackFlipDuration(t, spec.ImageData)
	}

	if angle := CleanFloat(params[rotate], 360); angle > 0 {
		t = time.Now()
		data = m.processor.Rotate(data, angle)
		m.metricService.TrackRotateDuration(t, spec.ImageData)
	}

	t = time.Now()
	src, err := m.processor.Encode(data, f)
	if err == nil {
		m.metricService.TrackEncodeDuration(t, spec.ImageData)
	}
	return src, err
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

func trackDuration(name string, start time.Time, spec processSpec) *metrics.UpdateOption {
	ext := strings.Split(http.DetectContentType(spec.ImageData), "/")[1]
	updateOption := metrics.UpdateOption{
		Name:     fmt.Sprintf("%s.%s.%s", name, metrics.GetImageSizeCluster(spec.ImageData), ext),
		Type:     metrics.Duration,
		Duration: time.Since(start),
		Scope:    spec.Scope,
	}
	metrics.Update(updateOption)
	return &updateOption
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
