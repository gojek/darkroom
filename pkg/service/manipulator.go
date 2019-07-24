package service

import (
	"bytes"
	"fmt"
	"github.com/gojek/darkroom/pkg/processor/native"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gojek/darkroom/pkg/metrics"
	"github.com/gojek/darkroom/pkg/processor"
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

	compress = "compress"

	cropDurationKey      = "cropDuration"
	decodeDurationKey    = "decodeDuration"
	encodeDurationKey    = "encodeDuration"
	grayScaleDurationKey = "grayScaleDuration"
	resizeDurationKey    = "resizeDuration"
	flipDurationKey      = "flipDuration"
	rotateDurationKey    = "rotateDuration"
	fixOrientationKey = "fixOrientation"
)

// Manipulator interface sets the contract on the implementation for common processing support in darkroom
type Manipulator interface {
	// Process takes ProcessSpec as an argument and returns []byte, error
	Process(spec ProcessSpec) ([]byte, error)
}

type manipulator struct {
	processor processor.Processor
}

// ProcessSpec defines the specification for a image manipulation job
type ProcessSpec struct {
	// Scope defines a scope for the image manipulation job, it can be used for logging/mertrics collection purposes
	Scope string
	// ImageData holds the actual image contents to processed
	ImageData []byte
	// Params hold the key-value pairs for the processing job and tells the manipulator what to do with the image
	Params map[string]string
}

// Process takes ProcessSpec as an argument and returns []byte, error
// This manipulator uses bild to do the actual image manipulations
func (m *manipulator) Process(spec ProcessSpec) ([]byte, error) {
	params := spec.Params
	var err error
	t := time.Now()
	data, f, err := m.processor.Decode(spec.ImageData)
	orientation, _ := native.GetOrientation(bytes.NewReader(spec.ImageData))
	if err != nil {
		return nil, err
	}
	trackDuration(decodeDurationKey, t, spec)
	if params[fit] == crop {
		t = time.Now()
		data = m.processor.Crop(data, CleanInt(params[width]), CleanInt(params[height]), GetCropPoint(params[crop]))
		trackDuration(cropDurationKey, t, spec)
	} else if len(params[fit]) == 0 && (CleanInt(params[width]) != 0 || CleanInt(params[height]) != 0) {
		t = time.Now()
		data = m.processor.Resize(data, CleanInt(params[width]), CleanInt(params[height]))
		trackDuration(resizeDurationKey, t, spec)
	}
	if params[mono] == blackHexCode {
		t = time.Now()
		data = m.processor.GrayScale(data)
		trackDuration(grayScaleDurationKey, t, spec)
	}

	if params[auto] == compress {
		t = time.Now()
		data = m.processor.FixOrientation(data, orientation)
		trackDuration(fixOrientationKey, t, spec)
	}

	if len(params[flip]) != 0 {
		t = time.Now()
		data = m.processor.Flip(data, params[flip])
		trackDuration(flipDurationKey, t, spec)
	}
	if CleanAngle(params[rotate]) > 0 {
		t = time.Now()
		data = m.processor.Rotate(data, CleanAngle(params[rotate]))
		trackDuration(rotateDurationKey, t, spec)
	}
	t = time.Now()
	src, err := m.processor.Encode(data, f)
	if err == nil {
		trackDuration(encodeDurationKey, t, spec)
	}
	return src, err
}

// CleanInt takes a string and return an int not greater than 9999
func CleanInt(input string) int {
	val, _ := strconv.Atoi(input)
	if val <= 0 {
		return 0
	}
	return val % 10000 // Never return value greater than 9999
}

// CleanAngle takes a string and return a float64 not greater than 360
func CleanAngle(input string) float64 {
	val, _ := strconv.ParseFloat(input, 64)
	if val <= 0 {
		return 0
	}
	return math.Mod(val, 360) // Never return value greater than 360
}

// GetCropPoint takes a string and returns the type CropPoint
func GetCropPoint(input string) processor.CropPoint {
	switch input {
	case "top":
		return processor.CropTop
	case "top,left":
		return processor.CropTopLeft
	case "top,right":
		return processor.CropTopRight
	case "left":
		return processor.CropLeft
	case "right":
		return processor.CropRight
	case "bottom":
		return processor.CropBottom
	case "bottom,left":
		return processor.CropBottomLeft
	case "bottom,right":
		return processor.CropBottomRight
	default:
		return processor.CropCenter
	}
}

func trackDuration(name string, start time.Time, spec ProcessSpec) *metrics.UpdateOption {
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

// NewManipulator takes in a Processor interface and returns a new manipulator
func NewManipulator(processor processor.Processor) *manipulator {
	return &manipulator{processor: processor}
}
