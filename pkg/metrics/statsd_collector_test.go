package metrics

import (
	"testing"
	"time"

	"github.com/cactus/go-statsd-client/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInitializeStatsdCollector(t *testing.T) {
	// Test Statter client
	_, err := InitializeStatsdCollector(&StatsdCollectorConfig{FlushBytes: 0})
	assert.Nil(t, err)
	assert.NotNil(t, instance)
	assert.NotNil(t, instance.client)

	// Test sampleRate
	_, err = InitializeStatsdCollector(&StatsdCollectorConfig{SampleRate: 5})
	assert.Nil(t, err)
	assert.Equal(t, float32(5), instance.sampleRate)

	_, err = InitializeStatsdCollector(&StatsdCollectorConfig{})
	assert.Nil(t, err)
	assert.Equal(t, float32(1), instance.sampleRate)
}

func TestRegisterHystrixMetrics(t *testing.T) {
	err := RegisterHystrixMetrics(&StatsdCollectorConfig{}, "prefix")
	assert.Nil(t, err)

	err = RegisterHystrixMetrics(&StatsdCollectorConfig{
		StatsdAddr: "foo:bar:foo",
	}, "prefix")
	assert.NotNil(t, err)
}

func TestStatsDMetricsUpdate(t *testing.T) {
	_,_ = InitializeStatsdCollector(&StatsdCollectorConfig{})

	mc := &mockStatsdClient{}
	instance.client = mc

	now := time.Now()
	mc.On("TimingDuration",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("time.Duration"),
		mock.AnythingOfType("float32")).Return(nil)
	instance.TrackDuration("cropDuration", now, nil)

	mc.On("Inc",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int64"),
		mock.AnythingOfType("float32")).Return(nil)
	instance.CountImageHandlerErrors("")

	mc.AssertExpectations(t)
}


type mockStatsdClient struct {
	mock.Mock
}

func (msc *mockStatsdClient) Inc(str string, i int64, sr float32) error {
	args := msc.Called(str, i, sr)
	return args.Error(0)
}

func (msc *mockStatsdClient) Dec(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Gauge(str string, i int64, sr float32) error {
	args := msc.Called(str, i, sr)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return args.Error(0)
}

func (msc *mockStatsdClient) GaugeDelta(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) Timing(string, int64, float32) error {
	panic("implement me")
}

func (msc *mockStatsdClient) TimingDuration(str string, t time.Duration, sr float32) error {
	args := msc.Called(str, t, sr)
	return args.Error(0)
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
