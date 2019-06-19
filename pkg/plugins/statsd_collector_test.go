package plugins

import (
	"github.com/cactus/go-statsd-client/statsd"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestInitializeStatsdCollector(t *testing.T) {
	scc, err := InitializeStatsdCollector(&StatsdCollectorConfig{})
	assert.Nil(t, err)
	assert.NotNil(t, scc)

	// Test sampleRate
	scc, err = InitializeStatsdCollector(&StatsdCollectorConfig{SampleRate: 5})
	assert.Nil(t, err)
	assert.NotNil(t, scc)
	assert.Equal(t, float32(5), scc.sampleRate)

	scc, err = InitializeStatsdCollector(&StatsdCollectorConfig{})
	assert.Nil(t, err)
	assert.NotNil(t, scc)
	assert.Equal(t, float32(1), scc.sampleRate)

	// Test Statter client
	scc, err = InitializeStatsdCollector(&StatsdCollectorConfig{})
	assert.Nil(t, err)
	assert.NotNil(t, scc)
	assert.NotNil(t, scc.client)
}

type mockStatsdClient struct {
}

func (msc *mockStatsdClient) Inc(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Dec(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Gauge(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) GaugeDelta(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Timing(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) TimingDuration(string, time.Duration, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Set(string, string, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) SetInt(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Raw(string, string, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) NewSubStatter(string) statsd.SubStatter {
	panic("implement me")
}

func (msc *mockStatsdClient) SetPrefix(string) {
	panic("implement me")
}

func (msc *mockStatsdClient) Close() error {
	panic("implement me")
}
