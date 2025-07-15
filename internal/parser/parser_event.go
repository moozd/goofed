package parser

import (
	"fmt"
	"strconv"
)

type ParserEvent struct {
	name          string
	expr          []byte
	params        []byte
	intermediates []byte
	char          byte
	final         byte
}

func newParserEvent() *ParserEvent {
	e := &ParserEvent{}
	e.rest()
	return e
}

func (e *ParserEvent) rest() {
	e.name = "unknown"
	e.char = 0x0
	e.expr = make([]byte, 0)
	e.clear()
}

func (e *ParserEvent) clear() {
	e.final = 0x0
	e.params = make([]byte, 0)
	e.intermediates = make([]byte, 0)
}

func (t *ParserEvent) String() string {
	return fmt.Sprintf("] %-12s: v=%-5s  F=%-5s P=%v I=%v", t.name, strconv.Quote(string(t.char)), strconv.Quote(string(t.final)), t.params, t.intermediates)
}
