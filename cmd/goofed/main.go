package main

import (
	"context"
	"log"
	"os"
	"runtime"

	"github.com/moozd/goofed/internal/gpu"
	"github.com/moozd/goofed/internal/screen"
	"github.com/moozd/goofed/internal/session"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	ctx := context.Background()

	shell, err := session.New(ctx, "zsh")
	defer shell.Close()

	if err != nil {
		log.Panicln("Could not start the PTY session.")
	}

	scr := screen.New(ctx, shell)
	defer scr.Close()

	cfg := gpu.Config{
		WindowTitle:  "Goofed",
		WindowHeight: 500,
		WindowWidth:  600,
	}

	g := gpu.New(ctx, cfg)

	g.Add(scr)

	if err = g.Run(); err != nil {
		log.Panicln(err)
	}

	os.Exit(0)
}
