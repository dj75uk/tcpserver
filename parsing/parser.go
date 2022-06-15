package parsing

import (
	"errors"
	"fmt"
	"strconv"
)

const stateBuildingCommand int = 0
const stateBuildingArg1LengthLength int = 1
const stateBuildingArg1Length int = 2
const stateBuildingArg1 int = 3
const stateBuildingArg2LengthLength int = 4
const stateBuildingArg2Length int = 5
const stateBuildingArg2 int = 6
const stateWaitingForMessageDequeue int = 7
const stateReset int = stateBuildingCommand

var ErrParserInvalidArgument = errors.New("invalid argument")
var ErrParserUnknownCommand = errors.New("unknown command")
var ErrParserBadFormat = errors.New("bad format")
var ErrParserNoMessage = errors.New("no message")

type ParserGrammar struct {
	ExpectedArguments uint16
	Arg1LengthIsValue bool
	Arg2LengthIsValue bool
}

type Parser struct {
	state             int
	command           string
	argsExpected      uint16
	arg1LengthIsValue bool
	arg1LengthLength  int
	arg1LengthBuilder string
	arg1Length        int
	arg1              string
	arg2LengthIsValue bool
	arg2LengthLength  int
	arg2LengthBuilder string
	arg2Length        int
	arg2              string
	commands          map[string]ParserGrammar
}

func NewParser(grammar map[string]ParserGrammar) (*Parser, error) {
	if grammar == nil {
		return nil, ErrParserInvalidArgument
	}
	result := &Parser{commands: grammar}
	result.reset()
	return result, nil
}

func CreateData(command string, key string, value string) ([]byte, error) {
	if len(command) != 3 {
		return nil, errors.New("invalid argument: 'command' must have length of 3")
	}
	if len(value) > 0 && len(key) == 0 {
		return nil, errors.New("invalid argument: cannot specify a 'value' with no 'key'")
	}
	result := command
	if keyLength := len(key); keyLength > 0 {
		keyLengthAsString := fmt.Sprintf("%d", keyLength)
		result += fmt.Sprintf("%d%s%s", len(keyLengthAsString), keyLengthAsString, key)
		if valueLength := len(value); valueLength > 0 {
			valueLengthAsString := fmt.Sprintf("%d", valueLength)
			result += fmt.Sprintf("%d%s%s", len(valueLengthAsString), valueLengthAsString, value)
		}
	}
	return []byte(result), nil
}

func (p *Parser) reset() {
	p.state = stateReset
	p.command = ""
	p.argsExpected = 0
	p.arg1LengthIsValue = false
	p.arg1LengthLength = 0
	p.arg1LengthBuilder = ""
	p.arg1Length = 0
	p.arg1 = ""
	p.arg2LengthIsValue = false
	p.arg2LengthLength = 0
	p.arg2LengthBuilder = ""
	p.arg2Length = 0
	p.arg2 = ""
}

func (p *Parser) GetMessage() (command string, arg1 string, arg2 string, err error) {
	if p.state == stateWaitingForMessageDequeue {
		defer p.reset()
		return p.command, p.arg1, p.arg2, nil
	}
	return "", "", "", ErrParserNoMessage
}

func (p *Parser) Process(datum string) (found bool, e error) {
	switch p.state {
	case stateBuildingCommand: // we're still waiting for a command...
		p.command += datum
		if len(p.command) == 3 {
			// validate command...
			if commandGrammar, exists := p.commands[p.command]; exists {
				if commandGrammar.ExpectedArguments == 0 {
					p.state = stateWaitingForMessageDequeue
					return true, nil // we have a valid zero-arg message
				}
				p.argsExpected = commandGrammar.ExpectedArguments
				p.arg1LengthIsValue = commandGrammar.Arg1LengthIsValue
				p.arg2LengthIsValue = commandGrammar.Arg2LengthIsValue
				p.state++
			} else {
				p.reset()
				return false, ErrParserUnknownCommand
			}
		}
	case stateBuildingArg1LengthLength: // we're waiting for the length of the arg1 length...
		if v, err := strconv.Atoi(datum); err == nil && v > 0 {
			p.arg1LengthLength = v
			p.state++
		} else {
			p.reset()
			return false, ErrParserBadFormat
		}
	case stateBuildingArg1Length: // we're waiting for the bytes of arg1 length...
		p.arg1LengthBuilder += datum
		if len(p.arg1LengthBuilder) == p.arg1LengthLength {
			v, err := strconv.Atoi(p.arg1LengthBuilder)

			if err == nil && !p.arg1LengthIsValue && v > 0 {
				p.arg1Length = v
				p.state++
			} else if err == nil && p.arg1LengthIsValue {
				p.arg1 = p.arg1LengthBuilder
				if p.argsExpected == 1 {
					p.state = stateWaitingForMessageDequeue
					return true, nil // we have a valid one-arg message in the special extension format
				}
				p.state = stateBuildingArg2LengthLength // skip to building arg2
			} else {
				p.reset()
				return false, ErrParserBadFormat
			}
		}
	case stateBuildingArg1: // we're waiting for the bytes of arg1...
		p.arg1 += datum
		if len(p.arg1) == p.arg1Length {
			if p.argsExpected == 1 {
				p.state = stateWaitingForMessageDequeue
				return true, nil // we have a valid one-arg message
			}
			p.state++
		}
	case stateBuildingArg2LengthLength: // we're waiting for the length of the arg2 length...
		if v, err := strconv.Atoi(datum); err == nil && v > 0 {
			p.arg2LengthLength = v
			p.state++
		} else {
			p.reset()
			return false, ErrParserBadFormat
		}
	case stateBuildingArg2Length: // we're waiting for the bytes of arg2 length...
		p.arg2LengthBuilder += datum
		if len(p.arg2LengthBuilder) == p.arg2LengthLength {
			v, err := strconv.Atoi(p.arg2LengthBuilder)

			if err == nil && !p.arg2LengthIsValue && v > 0 {
				p.arg2Length = v
				p.state++
			} else if err == nil && p.arg2LengthIsValue {
				p.arg2 = p.arg2LengthBuilder
				p.state = stateWaitingForMessageDequeue
				return true, nil // we have a valid two-arg message in the special extension format
			} else {
				p.reset()
				return false, ErrParserBadFormat
			}
		}
	case stateBuildingArg2: // we're waiting for the bytes of arg2...
		p.arg2 += datum
		if len(p.arg2) == p.arg2Length {
			p.state++
			if p.argsExpected == 2 {
				p.state = stateWaitingForMessageDequeue
				return true, nil // we have a valid two-arg message
			}
		}
	case stateWaitingForMessageDequeue: // we're waiting for GetMessage() to be called...
		// nop
	}
	return false, nil // we need more data
}
