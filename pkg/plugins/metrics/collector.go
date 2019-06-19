package metrics

import "time"

type MetricResult struct {
	ProcessingDuration time.Duration
	DownloadDuration   time.Duration
	TotalDuration      time.Duration
	ReadDuration       time.Duration
	WriteDuration      time.Duration
	CropDuration       time.Duration
	ResizeDuration     time.Duration
	MonoDuration       time.Duration
}

// MetricCollector represents the contract that all collectors must fulfill to gather statistics.
// Implementations of this interface do not have to maintain locking around their data stores so long as
// they are not modified outside of the storage/processor context.
type MetricCollector interface {
	// Update accepts a set of metrics from a command execution for remote instrumentation
	Update(MetricResult)
	// Reset resets the internal counters and timers.
	Reset()
}
