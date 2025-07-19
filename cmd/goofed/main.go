package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/moozd/goofed/internal/screen"
	"github.com/moozd/goofed/internal/session"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	ctx := context.Background()

	shell, err := session.New(ctx, "zsh")

	if err != nil {
		log.Panicln("Could not start the PTY session.")
	}
	defer shell.Close()

	scr := screen.New(ctx, shell)
	defer scr.Close()

	scr.Loop()

	os.Exit(0)
}
