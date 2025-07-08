package vte

import (
	"fmt"
	"strconv"
)

type Event struct {
	name          string
	expr          []rune
	params        []byte
	intermediates []byte
	char          byte
	final         byte
}

func NewEvent() *Event {
	e := &Event{}
	e.rest()
	return e
}

func (e *Event) rest() {
	e.name = "unknown"
	e.char = 0x0
	e.final = 0x0
	e.expr = make([]rune, 0)
	e.params = make([]byte, 0)
	e.intermediates = make([]byte, 0)
}

func (t *Event) String() string {
	return fmt.Sprintf("] %-12s: v=%-5s  F=%-5s P=%v I=%v", t.name, strconv.Quote(string(t.char)), strconv.Quote(string(t.final)), t.params, t.intermediates)
}
