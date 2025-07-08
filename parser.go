package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
	"time"
)

type State int
type Action int

const (
	StateGround State = iota
	StateAnyWhere
	StateUtf8
	StateSosPmApcString
	StateEscape
	StateEscapeIntermediate
	StateCsiEntry
	StateCsiIgnore
	StateCsiParam
	StateCsiIntermediate
	StateCsiFinal
	StateOscString
	StateDcsEntry
	StateDcsIgnore
	StateDcsIntermediate
	StateDcsParam
	StateDcsPassthrough
)
const (
	ActionNone Action = iota
	ActionClear
	ActionCollect
	ActionCsiDispatch
	ActionEscDispatch
	ActionExecute
	ActionHook
	ActionIgnore
	ActionOscEnd
	ActionOscPut
	ActionOscStart
	ActionParam
	ActionPrint
	ActionPut
	ActionUnhook
	ActionBeginUtf8
)

var states = map[State]string{
	StateGround:             "ground",
	StateAnyWhere:           "anywhere",
	StateUtf8:               "utf8",
	StateSosPmApcString:     "sos-pm-apc-string",
	StateEscape:             "escape",
	StateEscapeIntermediate: "escape-intermediate",
	StateCsiEntry:           "csi-entry",
	StateCsiIgnore:          "csi-ignore",
	StateCsiParam:           "csi-param",
	StateCsiIntermediate:    "csi-intermediate",
	StateCsiFinal:           "csi-final",
	StateOscString:          "osc-string",
	StateDcsEntry:           "dcs-entry",
	StateDcsIgnore:          "dcs-ignore",
	StateDcsIntermediate:    "dcs-intermediate",
	StateDcsParam:           "dcs-param",
	StateDcsPassthrough:     "dcs-passthrough",
}

var actions = map[Action]string{
	ActionNone:        "none",
	ActionClear:       "clear",
	ActionCollect:     "collect",
	ActionCsiDispatch: "csi-dispatch",
	ActionEscDispatch: "esc-dispatch",
	ActionExecute:     "execute",
	ActionHook:        "hook",
	ActionIgnore:      "ignore",
	ActionOscEnd:      "osc-end",
	ActionOscPut:      "osc-put",
	ActionOscStart:    "osc-start",
	ActionParam:       "param",
	ActionPrint:       "print",
	ActionPut:         "put",
	ActionUnhook:      "unhook",
	ActionBeginUtf8:   "begin-utf8",
}

type Parser struct {
	src      io.Reader
	state    State
	command  *Command
	Queue    chan Command
	ctx      context.Context
	cancelFn context.CancelFunc
}

func NewParser(ctx context.Context, src io.Reader) *Parser {
	ctx, cancel := context.WithCancel(ctx)

	parser := &Parser{
		src:      src,
		state:    StateGround,
		command:  NewCommand(),
		Queue:    make(chan Command, 200),
		ctx:      ctx,
		cancelFn: cancel,
	}

	go parser.worker()

	return parser
}

func (p *Parser) Close() {
	p.cancelFn()
	close(p.Queue)

}

func (p *Parser) worker() {
	defer close(p.Queue)
	reader := bufio.NewReader(p.src)
	for {
		select {
		case <-p.ctx.Done():
			return
		default:
			b, err := reader.ReadByte()
			if err == io.EOF {
				return
			}
			p.feed(b)
		}
	}
}

func (p *Parser) feed(c byte) {
	p.command.raw = append(p.command.raw, rune(c))
	state, action := p.advance(c)

	fmt.Printf("%-10s -- 0x%-2X -->  %-10s   %-15s  %s  \n", states[p.state], c, states[state], actions[action], strconv.QuoteRune(rune(c)))
	p.state = state

	switch action {
	case ActionPrint:
		p.send("print")
	case ActionCsiDispatch:
		p.send("csi")
	case ActionParam:
		p.command.params += string(c)
	}
	time.Sleep(1 * time.Second)

}

func (p *Parser) send(name string) {
	p.command.name = name
	// p.Queue <- *p.command
	p.command = NewCommand()
}

func (p *Parser) advance(c byte) (State, Action) {
	switch p.state {

	case StateAnyWhere:
		switch c {
		case 0x18, 0x1a:
			return StateGround, ActionExecute
		case 0x1b:
			return StateEscape, ActionNone
		}

	case StateGround:
		switch {
		case
			c == 0x19 || c == 0x9c,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x80, 0x8f),
			isBetween(c, 0x91, 0x9a):
			return StateAnyWhere, ActionExecute
		case isBetween(c, 0x20, 0x7f):
			return StateAnyWhere, ActionPrint
		case
			isBetween(c, 0xc2, 0xdf),
			isBetween(c, 0xe0, 0xef),
			isBetween(c, 0xf0, 0xf4):
			return StateUtf8, ActionBeginUtf8
		}

	case StateEscape:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionExecute
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x20, 0x2f):
			return StateEscapeIntermediate, ActionCollect
		case
			c == 0x59, c == 0x5a, c == 0x5c,
			isBetween(c, 0x30, 0x4f),
			isBetween(c, 0x51, 0x57),
			isBetween(c, 0x60, 0x7e):
			return StateGround, ActionEscDispatch
		case c == 0x5b:
			return StateCsiEntry, ActionNone
		case c == 0x5d:
			return StateOscString, ActionNone
		case c == 0x50:
			return StateDcsEntry, ActionNone
		case c == 0x58, c == 0x5e, c == 0x5f:
			return StateSosPmApcString, ActionNone
		}

	case StateEscapeIntermediate:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionExecute
		case isBetween(c, 0x20, 0x2f):
			return StateEscapeIntermediate, ActionCollect
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x30, 0x7e):
			return StateGround, ActionEscDispatch
		}

	case StateCsiEntry:
		switch {
		case isBetween(c, 0x00, 0x17), c == 0x19, isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionExecute
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x20, 0x2f):
			return StateCsiIntermediate, ActionCollect
		case isBetween(c, 0x30, 0x39),
			c == 0x3a, c == 0x3b:
			return StateCsiParam, ActionParam
		case isBetween(c, 0x3c, 0x3f):
			return StateCsiParam, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionCsiDispatch
		}

	case StateCsiIgnore:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionExecute
		case
			c == 0x7f,
			isBetween(c, 0x20, 0x3f):
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionNone
		}

	case StateCsiParam:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionExecute
		case isBetween(c, 0x30, 0x39), c == 0x3a, c == 0x3b:
			return StateAnyWhere, ActionParam
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x3c, 0x3f):
			return StateCsiIgnore, ActionNone
		case isBetween(c, 0x20, 0x2f):
			return StateCsiIntermediate, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionCsiDispatch
		}

	case StateCsiIntermediate:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionExecute
		case isBetween(c, 0x20, 0x2f):
			return StateAnyWhere, ActionCollect
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x30, 0x3f):
			return StateCsiIgnore, ActionNone
		case isBetween(c, 0x40, 0x7e):
			return StateGround, ActionCsiDispatch
		}

	case StateDcsEntry:
		switch {
		case
			c == 0x19,
			c == 0x7f,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x20, 0x2f):
			return StateDcsIntermediate, ActionCollect
		case
			c == 0x3a,
			c == 0x3b,
			isBetween(c, 0x30, 0x39):
			return StateDcsParam, ActionParam
		case isBetween(c, 0x3c, 0x3f):
			return StateDcsParam, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateDcsPassthrough, ActionNone
		}

	case StateDcsIntermediate:
		switch {
		case
			c == 0x7f,
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x20, 0x2f):
			return StateAnyWhere, ActionCollect
		case isBetween(c, 0x30, 0x3f):
			return StateDcsIgnore, ActionNone
		case isBetween(c, 0x40, 0x7e):
			return StateDcsPassthrough, ActionNone
		}

	case StateDcsIgnore:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x20, 0x7f):
			return StateAnyWhere, ActionIgnore
		case c == 0x9c:
			return StateGround, ActionNone
		}

	case StateDcsParam:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f):
			return StateAnyWhere, ActionIgnore
		case
			c == 0x3a,
			c == 0x3b,
			isBetween(c, 0x30, 0x39):
			return StateAnyWhere, ActionParam
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case isBetween(c, 0x3c, 0x3f):
			return StateDcsIgnore, ActionNone
		case isBetween(c, 0x20, 0x2f):
			return StateDcsIntermediate, ActionCollect
		case isBetween(c, 0x40, 0x7e):
			return StateDcsPassthrough, ActionNone
		}

	case StateDcsPassthrough:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x20, 0x7e):
			return StateAnyWhere, ActionPut
		case c == 0x7f:
			return StateAnyWhere, ActionIgnore
		case c == 0x9c:
			return StateGround, ActionNone
		}

	case StateSosPmApcString:
		switch {
		case
			c == 0x19,
			isBetween(c, 0x00, 0x17),
			isBetween(c, 0x1c, 0x1f),
			isBetween(c, 0x20, 0x7f):
			return StateAnyWhere, ActionIgnore
		case c == 0x9c:
			return StateGround, ActionNone
		}
	}

	return p.state, ActionNone
}
func isBetween(n, l, h byte) bool {
	return n > l && n < h
}

type Command struct {
	name          string
	raw           []rune
	params        string
	intermediates string
	final         string
}

func NewCommand() *Command {
	return &Command{
		name:          "",
		raw:           make([]rune, 0),
		params:        "",
		intermediates: "",
	}
}

func (t *Command) String() string {
	return fmt.Sprintf("[%-5s]: %-15s p=%s", t.name, strconv.Quote(string(t.raw)), t.params)
}
