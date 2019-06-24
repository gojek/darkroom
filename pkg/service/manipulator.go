package service

import (
	"***REMOVED***/darkroom/core/pkg/metrics"
	"***REMOVED***/darkroom/core/pkg/processor"
	"strconv"
	"time"
)

const (
	width        = "w"
	height       = "h"
	fit          = "fit"
	crop         = "crop"
	mono         = "mono"
	blackHexCode = "000000"

	cropDurationKey      = "cropDuration"
	resizeDurationKey    = "resizeDuration"
	watermarkDurationKey = "watermarkDuration"
	grayScaleDurationKey = "grayScaleDuration"
	decodeDurationKey    = "decodeDuration"
	encodeDurationKey    = "encodeDuration"
)

type Manipulator interface {
	Process(spec ProcessSpec) ([]byte, error)
}

type manipulator struct {
	processor processor.Processor
}

type ProcessSpec struct {
	Scope     string
	ImageData []byte
	Params    map[string]string
}

func (m *manipulator) Process(spec ProcessSpec) ([]byte, error) {
	params := spec.Params
	data := spec.ImageData
	var err error
	if params[fit] == crop {
		t := time.Now()
		data, err = m.processor.Crop(data, CleanInt(params[width]), CleanInt(params[height]), GetCropPoint(params[crop]))
		if err == nil {
			metrics.Update(metrics.UpdateOption{Name: cropDurationKey, Type: metrics.Duration, Duration: time.Since(t), Scope: spec.Scope})
		}
	} else if len(params[fit]) == 0 && (CleanInt(params[width]) != 0 || CleanInt(params[height]) != 0) {
		t := time.Now()
		data, err = m.processor.Resize(data, CleanInt(params[width]), CleanInt(params[height]))
		if err == nil {
			metrics.Update(metrics.UpdateOption{Name: resizeDurationKey, Type: metrics.Duration, Duration: time.Since(t), Scope: spec.Scope})
		}
	}
	if params[mono] == blackHexCode {
		t := time.Now()
		data, err = m.processor.GrayScale(data)
		if err == nil {
			metrics.Update(metrics.UpdateOption{Name: grayScaleDurationKey, Type: metrics.Duration, Duration: time.Since(t), Scope: spec.Scope})
		}
	}
	return data, err
}

func CleanInt(input string) int {
	val, _ := strconv.Atoi(input)
	if val <= 0 {
		return 0
	}
	return val % 10000 // Never return value greater than 9999
}

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

func NewManipulator(processor processor.Processor) *manipulator {
	return &manipulator{processor: processor}
}
