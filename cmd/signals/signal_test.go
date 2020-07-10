package signals

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/signal"
	"sync"
	"testing"
	"time"
)

func TestSetupSignalHandler(t *testing.T) {
	stop := SetupSignalHandler()
	task := &Task{
		ticker: time.NewTicker(time.Second * 2),
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	task.wg.Add(1)
	go func(c chan os.Signal) {
		defer task.wg.Done()
		task.Run(c)
	}(c)

	select {
	case sig := <-c:
		fmt.Printf("Got %s signal. Aborting...\n", sig)
	case _, ok := <-stop:
		assert.False(t, ok)
	}
}

type Task struct {
	wg     sync.WaitGroup
	ticker *time.Ticker
}

func (t *Task) Run(c chan os.Signal) {
	for {
		go sendSignal(c)
		handle()
	}
}

func handle() {
	for i := 0; i < 5; i++ {
		fmt.Print("#")
		time.Sleep(time.Millisecond * 100)
	}
	fmt.Println()
}

func sendSignal(stopChan chan os.Signal) {
	fmt.Printf("...")
	time.Sleep(1 * time.Second)
	stopChan <- os.Interrupt
}
