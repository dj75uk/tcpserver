package parsing

import (
	"errors"
	"strconv"
)

var ErrParserUnknownCommand = errors.New("unknown command")
var ErrParserBadFormat = errors.New("bad format")

type Parser2 struct {
	state             int
	command           string
	argsExpected      uint16
	arg1LengthLength  int
	arg1LengthBuilder string
	arg1Length        int
	arg1              string
	arg2LengthLength  int
	arg2LengthBuilder string
	arg2Length        int
	arg2              string
	commands          map[string]uint16
}

func NewParser2() *Parser2 {
	result := &Parser2{
		commands: map[string]uint16{
			"put": 2,
			"get": 1,
			"del": 1,
			"bye": 0,
		}}
	result.reset()
	return result
}

func (p *Parser2) reset() {
	p.state = 0
	p.command = ""
	p.argsExpected = 0
	p.arg1LengthLength = 0
	p.arg1LengthBuilder = ""
	p.arg1Length = 0
	p.arg1 = ""
	p.arg2LengthLength = 0
	p.arg2LengthBuilder = ""
	p.arg2Length = 0
	p.arg2 = ""
}

func (p *Parser2) GetMessage() (command string, arg1 string, arg2 string, err error) {
	defer p.reset()
	return p.command, p.arg1, p.arg2, nil
}

func (p *Parser2) Process(datum string) (found bool, e error) {

	switch p.state {

	case 0: // we're still waiting for a command...
		p.command += datum
		if len(p.command) == 3 {
			// validate command...
			if argsExpected, exists := p.commands[p.command]; exists {
				if argsExpected == 0 {
					return true, nil // we have a valid zero-arg message
				}
				p.argsExpected = argsExpected
				p.state++
			} else {
				p.reset()
				return false, ErrParserUnknownCommand
			}
		}
	case 1: // we're waiting for the length of the arg1 length...
		if v, err := strconv.Atoi(datum); err == nil && v > 0 {
			p.arg1LengthLength = v
			p.state++
		} else {
			p.reset()
			return false, ErrParserBadFormat
		}
	case 2: // we're waiting for the bytes of arg1 length...
		p.arg1LengthBuilder += datum
		if len(p.arg1LengthBuilder) == p.arg1LengthLength {
			if v, err := strconv.Atoi(p.arg1LengthBuilder); err == nil && v > 0 {
				p.arg1Length = v
				p.state++
			} else {
				p.reset()
				return false, ErrParserBadFormat
			}
		}
	case 3: // we're waiting for the bytes of arg1...
		p.arg1 += datum
		if len(p.arg1) == p.arg1Length {
			if p.argsExpected == 1 {
				return true, nil
			}
			p.state++
		}

	case 4: // we're waiting for the length of the arg2 length...
		if v, err := strconv.Atoi(datum); err == nil && v > 0 {
			p.arg2LengthLength = v
			p.state++
		} else {
			p.reset()
			return false, ErrParserBadFormat
		}
	case 5: // we're waiting for the bytes of arg2 length...
		p.arg2LengthBuilder += datum
		if len(p.arg2LengthBuilder) == p.arg2LengthLength {
			if v, err := strconv.Atoi(p.arg2LengthBuilder); err == nil && v > 0 {
				p.arg2Length = v
				p.state++
			} else {
				p.reset()
				return false, ErrParserBadFormat
			}
		}
	case 6: // we're waiting for the bytes of arg2...
		p.arg2 += datum
		if len(p.arg2) == p.arg2Length {
			if p.argsExpected == 2 {
				return true, nil
			}
			p.state++
		}

	}

	return false, nil
}
