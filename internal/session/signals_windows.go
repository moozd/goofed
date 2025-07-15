//go:build windows
// +build windows

package session

import (
	"os"
	"os/signal"
)

func getStopSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	return signals
}

func getResizeSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
	close(signals)
	return signals
}
