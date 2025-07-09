//go:build linux || darwin
// +build linux darwin

package pty

import (
	"os"
	"os/signal"
	"syscall"
)

func getStopSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	return signals
}

func getResizeSignals() chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGWINCH)
	return signals
}
