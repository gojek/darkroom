package metrics

import "time"

// Type is used to differentiate between the various metrics update options possible
type Type int

const (
	// Duration can update a key which requires a time.Duration value
	Duration Type = 0
	// Gauge can update a key which requires a numeric value
	Gauge    Type = 1
	// Count can update a key which requires simple increment operation for counting occurrences
	Count    Type = 2
)

// UpdateOption is used to specify the specs for a metrics update operation
type UpdateOption struct {
	// Scope can be used to provide extra context to the metrics
	Scope    string
	// Name is the actual key that will be pushed to stats
	// format: {StatsdCollectorConfig.Prefix}.{Scope|default}.{Name}
	Name     string
	// Type is used to differentiate between the various metrics update options possible
	Type     Type
	// NumValue holds a float64 value for the update option
	NumValue float64
	// Duration holds a time.Duration value for the update option
	Duration time.Duration
}
