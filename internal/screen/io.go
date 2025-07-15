package screen

import "github.com/moozd/goofed/internal/parser"

func (self *Screen) handle(event parser.ParserEvent) {

}

func (self *Screen) send(b ...byte) {
	self.session.Write(b)
}
