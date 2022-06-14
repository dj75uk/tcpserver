package parsing

import (
	"errors"
	"fmt"
	"strconv"
)

type Msg struct {
	Command string
	Key     string
	Value   string
}

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (parser *Parser) CreateData(command string, key string, value string) ([]byte, error) {
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

func (parser *Parser) Parse(data []byte, expectedParameters uint16) (message *Msg, bytesProcessed int, err error) {
	length := len(data)
	if length < 3 {
		return nil, 0, nil
	}
	command := string(data[0:3])

	if expectedParameters == 0 {
		return &Msg{
			Command: command,
			Key:     "",
			Value:   "",
		}, 3, nil
	}

	stage := 0
	kss := 0
	ksBuffer := ""
	keyLength := 0
	key := ""
	vss := 0
	vsBuffer := ""
	valueLength := 0
	value := ""

	done := false
	processed := 3
	for index := 3; index < length; index++ {

		if done {
			break
		}
		datum := string(data[index])
		processed++

		switch stage {
		case 0: // read kss
			kss, err = strconv.Atoi(datum)
			if err != nil || kss == 0 {
				return nil, processed, errors.New("bad format")
			}
			stage++
		case 1: // read ks based upon value of kss
			ksBuffer += datum
			if len(ksBuffer) == kss {
				keyLength, err = strconv.Atoi(ksBuffer)
				if err != nil || keyLength <= 0 {
					return nil, processed, errors.New("bad format")
				}
				stage++
			}
		case 2: // read k based upon value of ks
			key += datum
			if len(key) == keyLength {
				// got the key
				stage++
				if expectedParameters == 1 {
					done = true
				}
			}
		case 3: // read vss
			vss, err = strconv.Atoi(datum)
			if err != nil || vss == 0 {
				return nil, processed, errors.New("bad format")
			}
			stage++
		case 4: // read vs based upon value of vss
			vsBuffer += datum
			if len(vsBuffer) == vss {
				valueLength, err = strconv.Atoi(vsBuffer)
				if err != nil || valueLength <= 0 {
					return nil, processed, errors.New("bad format")
				}
				stage++
			}
		case 5: // read v based upon value of vs
			value += datum
			if len(value) == valueLength {
				// got the key
				done = true
			}
		default:
			done = true
		}

	}

	// return a completed message...
	if done {
		return &Msg{
			Command: command,
			Key:     key,
			Value:   value,
		}, processed, nil
	}

	// indicate that we need more data...
	return nil, 0, nil
}
