package parser

import (
	"bufio"
	"context"
	"io"
	"log"
)

type State string
type Action string

const (
	StateGround             State = "ground"
	StateSosPmApcString           = "sos-pm-apc"
	StateEscape                   = "escape"
	StateEscapeIntermediate       = "escape-intermediate"
	StateCsiEntry                 = "csi-entry"
	StateCsiIgnore                = "csi-ignore"
	StateCsiParam                 = "csi-param"
	StateCsiIntermediate          = "csi-intermediate"
	StateOscString                = "osc-string"
	StateDcsEntry                 = "dcs-entry"
	StateDcsIgnore                = "dcs-ignore"
	StateDcsIntermediate          = "dcs-intermediate"
	StateDcsParam                 = "dcs-param"
	StateDcsPassthrough           = "dcs-passthrough"
)
const (
	ActionNone        Action = "none"
	ActionClear              = "clear"
	ActionCollect            = "collect"
	ActionCsiDispatch        = "csi.dispatch"
	ActionEscDispatch        = "esc.dispatch"
	ActionExecute            = "execute"
	ActionHook               = "hook"
	ActionIgnore             = "ignore"
	ActionOscEnd             = "osc.end"
	ActionOscPut             = "osc.put"
	ActionOscStart           = "osc.start"
	ActionParam              = "param"
	ActionPrint              = "print"
	ActionPut                = "put"
	ActionUnhook             = "unhook"
)

type Parser struct {
	src    io.Reader
	state  State
	event  *ParserEvent
	ctx    context.Context
	cancel context.CancelFunc

	Queue chan ParserEvent
}

func New(ctx context.Context, src io.Reader) *Parser {
	ctx, cancel := context.WithCancel(ctx)

	self := &Parser{
		src:    src,
		state:  StateGround,
		event:  newParserEvent(),
		Queue:  make(chan ParserEvent, 256),
		ctx:    ctx,
		cancel: cancel,
	}

	go self.worker()

	return self
}

func (self *Parser) Close() {
	self.cancel()
	close(self.Queue)
}

func (self *Parser) worker() {
	reader := bufio.NewReader(self.src)
	for {
		select {
		case <-self.ctx.Done():
			return
		default:
			b, err := reader.ReadByte()
			if err == io.EOF {
				return
			}
			self.feed(b)
		}
	}
}

func (self *Parser) feed(c byte) {
	// transition to the next state by visiting the new char
	state, action := self.transition(c)

	// preform the action
	self.act(action, c)

	// change the state
	self.state = state
}

func (self *Parser) dispatch(action Action) {
	self.event.name = string(action)

	select {
	case <-self.ctx.Done():
		log.Fatal(self.ctx.Err())
	case self.Queue <- *self.event:
		self.event.rest()
	}
}

func (self *Parser) act(action Action, c byte) {
	self.event.char = c
	self.event.expr = append(self.event.expr, c)

	switch action {
	case ActionClear:
		self.event.clear()
	case ActionCollect:
		self.event.intermediates = append(self.event.intermediates, c)
	case ActionParam:
		self.event.params = append(self.event.params, c)
	case
		ActionCsiDispatch,
		ActionEscDispatch,
		ActionUnhook:
		self.event.final = c
		self.dispatch(action)
	case
		ActionPut,
		ActionPrint,
		ActionHook,
		ActionOscStart,
		ActionOscPut,
		ActionOscEnd,
		ActionExecute:
		self.dispatch(action)
	case
		ActionIgnore,
		ActionNone:
	}

}

func (self *Parser) transition(c byte) (State, Action) {
	switch self.state {

	case StateGround:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateGround, ActionExecute
		case isBetween(c, 0x20, 0x7f):
			return StateGround, ActionPrint
		}

	case StateEscape:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateEscape, ActionExecute
		case c == 0x7f:
			return StateEscape, ActionIgnore

		case isBetween(c, 0x20, 0x2f):
			return StateEscapeIntermediate, ActionCollect
		case c == 0x5b:
			return StateCsiEntry, ActionClear
		case c == 0x50:
			return StateDcsEntry, ActionClear
		case c == 0x5d:
			return StateOscString, ActionOscStart
		case c == 0x58, c == 0x5e, c == 0x5f:
			return StateSosPmApcString, ActionNone
		case
			c == 0x59,
			c == 0x5a,
			c == 0x5c,
			isBetween(c, 0x30, 0x4f),
			isBetween(c, 0x51, 0x57),
			isBetween(c, 0x60, 0x7e):
			return StateGround, ActionEscDispatch
		}

	case StateEscapeIntermediate:
		switch {
		case c == 0x7f:
			return StateEscapeIntermediate, ActionIgnore
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateEscapeIntermediate, ActionExecute
		case isBetween(c, 0x20, 0x2f):
			return StateEscapeIntermediate, ActionCollect

		case isBetween(c, 0x30, 0x7e):
			return StateGround, ActionEscDispatch
		}

	case StateCsiEntry:
		switch {
		case c == 0x7f:
			return StateCsiEntry, ActionIgnore
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateCsiEntry, ActionExecute

		case c == 0x3a:
			return StateCsiIgnore, ActionNone
		case isBetween(c, 0x20, 0x2f):
			return StateCsiIntermediate, ActionCollect
		case
			c == 0x3b,
			isBetween(c, 0x30, 0x39):
			return StateCsiParam, ActionParam
		case isBetween(c, 0x3c, 0x3f):
			return StateCsiParam, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionCsiDispatch
		}
	case StateCsiIgnore:
		switch {
		case
			c == 0x7f,
			isBetween(c, 0x20, 0x3f):
			return StateCsiIgnore, ActionIgnore
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateCsiIgnore, ActionExecute

		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionNone
		}
	case StateCsiIntermediate:
		switch {
		case c == 0x7f:
			return StateCsiIntermediate, ActionIgnore
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateCsiIntermediate, ActionExecute
		case isBetween(c, 0x20, 0x2f):
			return StateCsiIntermediate, ActionCollect

		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionCsiDispatch
		}
	case StateCsiParam:
		switch {

		case c == 0x7f:
			return StateCsiParam, ActionIgnore
		case
			c == 0x19,
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x00, 0x17):
			return StateCsiParam, ActionExecute
		case
			c == 0x3b,
			isBetween(c, 0x30, 0x39):
			return StateCsiParam, ActionParam

		case isBetween(c, 0x20, 0x2f):
			return StateCsiIntermediate, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionCsiDispatch
		}
	case StateOscString:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateOscString, ActionIgnore
		case isBetween(c, 0x20, 0x7f):
			return StateOscString, ActionOscPut

		case c == 0x9c:
			return StateGround, ActionOscEnd
		}
	case StateDcsEntry:
		switch {
		case
			c == 0x7f,
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateDcsEntry, ActionIgnore

		case c == 0x3a:
			return StateDcsIgnore, ActionIgnore
		case isBetween(c, 0x20, 0x2f):
			return StateDcsIntermediate, ActionCollect
		case
			c == 0x3b,
			isBetween(c, 0x30, 0x39):
			return StateDcsParam, ActionParam
		case isBetween(c, 0x3c, 0x3f):
			return StateDcsParam, ActionParam
		case isBetween(c, 0x40, 0x7e):
			return StateDcsPassthrough, ActionHook

		}
	case StateDcsIgnore:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x20, 0x7f),
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateDcsIgnore, ActionIgnore

		case c == 0x9c:
			return StateGround, ActionNone
		}

	case StateDcsParam:
		switch {
		case
			c == 0x7f, c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateDcsParam, ActionIgnore
		case
			c == 0x3b,
			isBetween(c, 0x30, 0x39):
			return StateDcsParam, ActionParam

		case
			c == 0x3a,
			isBetween(c, 0x3c, 0x3f):
			return StateDcsIgnore, ActionNone
		case isBetween(c, 0x20, 0x2f):
			return StateDcsIntermediate, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateDcsPassthrough, ActionHook
		}
	case StateDcsIntermediate:
		switch {
		case
			c == 0x7f, c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateCsiIntermediate, ActionIgnore
		case isBetween(c, 0x20, 0x2f):
			return StateDcsIntermediate, ActionCollect

		case isBetween(c, 0x30, 0x3f):
			return StateDcsIgnore, ActionNone
		case isBetween(c, 0x40, 0x7e):
			return StateDcsPassthrough, ActionNone
		}
	case StateDcsPassthrough:
		switch {
		case c == 0x7f:
			return StateDcsPassthrough, ActionIgnore
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x20, 0x7e),
			isBetween(c, 0x1c, 0x1f):
			return StateDcsPassthrough, ActionPut
		}
	case StateSosPmApcString:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x20, 0x7f):
			return StateOscString, ActionNone

		case c == 0x9c:
			return StateGround, ActionNone
		}
	}

	// anywhere
	switch {
	case
		c == 0x18, c == 0x1a,
		isBetween(c, 0x80, 0x8f),
		isBetween(c, 0x91, 0x97),
		isBetween(c, 0x99, 0x9A):
		return StateGround, ActionExecute
	case c == 0x1b:
		return StateEscape, ActionClear
	case c == 0x9b:
		return StateCsiEntry, ActionClear
	case c == 0x9d:
		return StateOscString, ActionNone
	case c == 0x98, c == 0x9e, c == 0x9f:
		return StateSosPmApcString, ActionNone
	case c == 0x90:
		return StateDcsEntry, ActionClear
	case c == 0x9c:
		return StateGround, ActionUnhook
	}

	return self.state, ActionNone
}
