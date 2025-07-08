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
	return &Event{
		name:          "",
		expr:          make([]rune, 0),
		params:        make([]byte, 0),
		intermediates: make([]byte, 0),
	}
}

func (t *Event) String() string {
	return fmt.Sprintf("] %-12s: F=%-5s P=%v I=%v", t.name, strconv.Quote(string(t.final)), t.params, t.intermediates)
}
