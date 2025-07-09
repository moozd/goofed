package graphics

import (
	"context"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/moozd/goofed/internal/pty"
	"github.com/moozd/goofed/internal/vte"
)

type Screen struct {
	height  int
	width   int
	ctx     context.Context
	session *pty.Session
	parser  *vte.Parser
}

func NewScreen(ctx context.Context, width, height int, session *pty.Session, parser *vte.Parser) *Screen {
	return &Screen{
		ctx:     ctx,
		width:   width,
		height:  height,
		session: session,
		parser:  parser,
	}
}

func (s *Screen) Update() error {
	return nil
}

func (s *Screen) Draw(g *ebiten.Image) {

}

func (s *Screen) Layout(ow, oh int) (int, int) {
	return ow, oh
}
