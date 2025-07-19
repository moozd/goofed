package screen

import (
	"context"

	"github.com/moozd/goofed/internal/parser"
	"github.com/moozd/goofed/internal/session"
)

type Screen struct {
	ctx     context.Context
	grid    *Grid
	parser  *parser.Parser
	session *session.Session
}

func New(c context.Context, s *session.Session) *Screen {

	self := &Screen{
		ctx:     c,
		session: s,
		grid:    newGrid(),
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
