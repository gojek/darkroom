package metrics

import "time"

type Type int

const (
	Duration Type = 0
	Gauge    Type = 1
	Count    Type = 2
)

type UpdateOption struct {
	Scope    string
	Name     string
	Type     Type
	NumValue float64
	Duration time.Duration
}
