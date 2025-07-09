package graphics

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/moozd/goofed/internal/pty"
	"github.com/moozd/goofed/internal/vte"
)

func MainLoop() error {

	ctx := context.Background()

	session, err := pty.NewSession(ctx, "zsh")
	defer session.Close()

	if err != nil {
		return err
	}

	parser := vte.NewParser(ctx, session)
	defer parser.Close()

	screen := NewScreen(ctx, 1024, 768, session, parser)

	ebiten.SetWindowTitle("goofed")
	ebiten.SetWindowSize(screen.width, screen.height)

	return ebiten.RunGame(screen)
}
