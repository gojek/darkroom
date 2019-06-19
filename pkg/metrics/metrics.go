package metrics

import "time"

type Type int

const (
	Duration Type = 0
	Guage    Type = 1
	Counter  Type = 2
)

type UpdateOption struct {
	Name     string
	Type     Type
	NumValue float64
	Duration time.Duration
}
