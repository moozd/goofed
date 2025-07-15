package screen

import (
	"context"

	"github.com/moozd/goofed/internal/gpu"
	"github.com/moozd/goofed/internal/parser"
	"github.com/moozd/goofed/internal/screen/grid"
	"github.com/moozd/goofed/internal/session"
)

type Screen struct {
	ctx        context.Context
	session    *session.Session
	parser     *parser.Parser
	grid       *grid.Grid
	gpuContext *gpuContext

	gpu.Model
}

func New(c context.Context, s *session.Session) *Screen {

	self := &Screen{
		ctx:     c,
		session: s,
		grid:    grid.New(),
		parser:  parser.New(c, s),
	}
	go self.drainParserQueue()

	return self
}

func (self *Screen) Close() {
	self.parser.Close()
}

func (self *Screen) drainParserQueue() {
	select {
	case <-self.ctx.Done():
		return
	default:
		for event := range self.parser.Queue {
			self.handle(event)
		}
	}
}
