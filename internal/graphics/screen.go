package graphics

import (
	"context"

	"github.com/moozd/goofed/internal/pty"
	"github.com/moozd/goofed/internal/vte"
)

type Screen struct {
	ctx     context.Context
	session *pty.Session
	parser  *vte.Parser
}

func NewScreen(c context.Context, p *vte.Parser, s *pty.Session) *Screen {
	return &Screen{
		ctx:     c,
		parser:  p,
		session: s,
	}
}

func (s *Screen) Init(r *Renderer) {}

func (s *Screen) Draw(r *Renderer) {}
