package toolbox

import (
	"os"
	"testing"
	"time"
)

func TestWaitForEnd(t *testing.T) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		if err := p.Signal(os.Interrupt); err != nil {
			panic("Could not signal interrupt")
		}
	}()
	sigch = make(chan os.Signal, 2)

	WaitForSignal()
}
