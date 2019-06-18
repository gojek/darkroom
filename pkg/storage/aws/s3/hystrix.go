package s3

import (
	"github.com/afex/hystrix-go/hystrix"
)

func makeNetworkCall(name string, commandConfig hystrix.CommandConfig, run func() error, fallback func(error) error) {
	hystrix.ConfigureCommand(name, commandConfig)
	hystrix.Go(name, run, fallback)
}
