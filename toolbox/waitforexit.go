package toolbox

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Make this available for tests
var sigch chan os.Signal

func init() {
	sigch = make(chan os.Signal, 2)
}

// WaitForSignal waits for a signal to terminate
func WaitForSignal() {
	logrus.Debug("Waiting for kill signal")
	terminator := make(chan bool)

	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigch
		terminator <- true
	}()
	<-terminator
}

// SendInterrupt sends an interrupt signal to the waiting channel
func SendInterrupt() {
	select {
	case sigch <- os.Interrupt:
		// ok
	default:
		// ignore
	}
}

// GetSignalChannel returns the signal channel. This is for testing.
func GetSignalChannel() chan os.Signal {
	return sigch
}
