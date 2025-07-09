package main

import (
	"runtime"

	"github.com/moozd/goofed/internal/app"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	app.MainLoop()
}
