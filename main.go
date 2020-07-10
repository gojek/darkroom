// Darkroom is an image proxy that works with different
// storage backends and different image processing engines. It also
// gives special attention to resiliency and speed. There is also
// support for inbuilt metrics collection for statsd.
package main

import (
	"github.com/gojek/darkroom/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
