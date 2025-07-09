package app

import (
	"context"
	"log"
	"os"

	"github.com/moozd/goofed/internal/graphics"
	"github.com/moozd/goofed/internal/pty"
	"github.com/moozd/goofed/internal/vte"
)

const (
	WIDTH  = 1024
	HEIGHT = 768
)

func MainLoop() {

	ctx := context.Background()

	session, err := pty.NewSession(ctx, "zsh")
	defer session.Close()

	if err != nil {
		log.Panicln("Could not start the PTY session.")

	}

	parser := vte.NewParser(ctx, session)
	defer parser.Close()

	screen := graphics.NewScreen(ctx, parser, session)

	renderer := graphics.NewRenderer(ctx, WIDTH, HEIGHT)
	renderer.Add(screen)

	renderer.Start()
	defer renderer.Close()

	os.Exit(0)
}
