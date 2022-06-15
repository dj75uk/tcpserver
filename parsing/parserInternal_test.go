package parsing

import (
	"fmt"
	"kvsapp/assertions"
	"testing"
)

func createTestObject() *Parser {
	result, _ := NewParser(map[string]uint16{
		"abc": 0,
		"def": 1,
		"ghi": 2,
	})
	return result
}

func assertGrammar(t *testing.T, testObject *Parser, testGrammar map[string]uint16) {
	if len(testObject.commands) != len(testGrammar) {
		t.Errorf("param: %s, expected: %d, actual: %d", "len(commands)", len(testGrammar), len(testObject.commands))
	}
	if len(testObject.commands) == len(testGrammar) {
		for expectedKey, expectedValue := range testGrammar {
			actualValue, exists := testObject.commands[expectedKey]
			if !exists {
				t.Errorf("param: commands[%s], expected: %s, actual: %s", expectedKey, expectedKey, "<key not in map>")
			}
			if exists && actualValue != expectedValue {
				t.Errorf("param: commands[%s], expected: %d, actual: %d", expectedKey, expectedValue, actualValue)
			}
		}
	}
}

func assertState(t *testing.T, testObject *Parser, expectedState int, expectedCommand string, expectedArgsExpected uint16) {
	assert := assertions.NewAssert(t)
	if testObject.state != expectedState {
		t.Errorf("param: %s, expected: %d, actual: %d", "state", expectedState, testObject.state)
	}
	assert.String("command", expectedCommand, testObject.command)
	if testObject.argsExpected != expectedArgsExpected {
		t.Errorf("param: %s, expected: %d, actual: %d", "argsExpected", expectedArgsExpected, testObject.argsExpected)
	}
}

func assertArg1(t *testing.T, testObject *Parser, expectedArg1LengthLength int, expectedArg1LengthBuilder string, expectedArg1Length int, expectedArg1 string) {
	assert := assertions.NewAssert(t)
	if testObject.arg1LengthLength != expectedArg1LengthLength {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg1LengthLength", expectedArg1LengthLength, testObject.arg1LengthLength)
	}
	if testObject.arg1LengthBuilder != expectedArg1LengthBuilder {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg1LengthBuilder", expectedArg1LengthBuilder, testObject.arg1LengthBuilder)
	}
	if testObject.arg1Length != expectedArg1Length {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg1Length", expectedArg1Length, testObject.arg1Length)
	}
	assert.String("arg1", expectedArg1, testObject.arg1)
}

func assertArg2(t *testing.T, testObject *Parser, expectedArg2LengthLength int, expectedArg2LengthBuilder string, expectedArg2Length int, expectedArg2 string) {
	assert := assertions.NewAssert(t)
	if testObject.arg2LengthLength != expectedArg2LengthLength {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg2LengthLength", expectedArg2LengthLength, testObject.arg2LengthLength)
	}
	if testObject.arg2LengthBuilder != expectedArg2LengthBuilder {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg2LengthBuilder", expectedArg2LengthBuilder, testObject.arg2LengthBuilder)
	}
	if testObject.arg2Length != expectedArg2Length {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg2Length", expectedArg2Length, testObject.arg2Length)
	}
	assert.String("arg2", expectedArg2, testObject.arg2)
}

func assertResetState(t *testing.T, testObject *Parser, testGrammar map[string]uint16) {
	assertState(t, testObject, stateReset, "", 0)
	assertArg1(t, testObject, 0, "", 0, "")
	assertArg2(t, testObject, 0, "", 0, "")
	assertGrammar(t, testObject, testGrammar)
}

func TestNewParserInitialisesStructure(t *testing.T) {
	t.Parallel()
	testGrammar := map[string]uint16{
		"aaa": 0,
		"bbb": 1,
		"ccc": 2,
	}

	testObject, _ := NewParser(testGrammar)
	assertResetState(t, testObject, testGrammar)
}

func TestResetClearsState(t *testing.T) {
	t.Parallel()
	testGrammar := map[string]uint16{
		"ddd": 0,
		"eee": 1,
		"fff": 2,
	}
	testObject, _ := NewParser(testGrammar)
	testObject.state = 123
	testObject.command = "abc"
	testObject.argsExpected = 234
	testObject.arg1LengthLength = 456
	testObject.arg1LengthBuilder = "def"
	testObject.arg1Length = 567
	testObject.arg1 = "ghi"
	testObject.arg2LengthLength = 678
	testObject.arg2LengthBuilder = "jkl"
	testObject.arg2Length = 789
	testObject.arg2 = "mno"
	testObject.reset()
	assertResetState(t, testObject, testGrammar)
}

func TestProcessUnknownCommand(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	testGrammar := map[string]uint16{
		"ddd": 0,
		"eee": 1,
		"fff": 2,
	}
	testObject, _ := NewParser(testGrammar)

	found := false
	err := error(nil)

	found, err = testObject.Process("x")
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, 0, "x", 0)

	found, err = testObject.Process("y")
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, 0, "xy", 0)

	found, err = testObject.Process("z")
	assert.False("found", found)
	assert.Error(ErrParserUnknownCommand, err)
	assertState(t, testObject, 0, "", 0)
}

func TestProcessKnownZeroArgumentCommand(t *testing.T) {
	assert := assertions.NewAssert(t)
	const ExpectedCommand string = "abc"
	t.Parallel()
	testGrammar := map[string]uint16{ExpectedCommand: 0}
	testObject, _ := NewParser(testGrammar)

	found := false
	err := error(nil)

	found, err = testObject.Process(ExpectedCommand[0:1])
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingCommand, ExpectedCommand[0:1], 0)

	found, err = testObject.Process(ExpectedCommand[1:2])
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingCommand, ExpectedCommand[0:2], 0)

	found, err = testObject.Process(ExpectedCommand[2:3])
	assert.True("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateWaitingForMessageDequeue, ExpectedCommand, 0)
}

func TestProcessKnownOneArgumentCommand(t *testing.T) {

	const ExpectedCommand string = "bcd"
	const ExpectedArg1 string = "arg1"

	assert := assertions.NewAssert(t)
	t.Parallel()
	testGrammar := map[string]uint16{ExpectedCommand: 1}
	testObject, _ := NewParser(testGrammar)

	found := false
	err := error(nil)

	found, err = testObject.Process(ExpectedCommand[0:1])
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingCommand, ExpectedCommand[0:1], 0)

	found, err = testObject.Process(ExpectedCommand[1:2])
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingCommand, ExpectedCommand[0:2], 0)

	found, err = testObject.Process(ExpectedCommand[2:3])
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingArg1LengthLength, ExpectedCommand, 1)

	found, err = testObject.Process("1")
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingArg1Length, ExpectedCommand, 1)

	found, err = testObject.Process(fmt.Sprintf("%d", len(ExpectedArg1)))
	assert.False("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateBuildingArg1, ExpectedCommand, 1)

	for _, r := range ExpectedArg1 {
		found, err = testObject.Process(string(r))
	}
	assert.True("found", found)
	assert.Error(nil, err)
	assertState(t, testObject, stateWaitingForMessageDequeue, ExpectedCommand, 1)
}

func TestGetMessageResetsState(t *testing.T) {
	const ExpectedCommand string = "qwe"
	const ExpectedArg1 string = "arg1value"
	const expectedArg2 string = "arg2value"
	t.Parallel()

	testGrammar := map[string]uint16{"ghj": 0}
	testObject, _ := NewParser(testGrammar)
	testObject.state = stateWaitingForMessageDequeue
	testObject.command = ExpectedCommand
	testObject.arg1 = ExpectedArg1
	testObject.arg2 = expectedArg2
	_, _, _, _ = testObject.GetMessage()

	assertResetState(t, testObject, testGrammar)
}

func TestGetMessageReturnsInternalState(t *testing.T) {
	const ExpectedCommand string = "qwe"
	const ExpectedArg1 string = "arg1value"
	const expectedArg2 string = "arg2value"
	t.Parallel()
	assert := assertions.NewAssert(t)

	testObject, _ := NewParser(map[string]uint16{})
	testObject.state = stateWaitingForMessageDequeue
	testObject.command = ExpectedCommand
	testObject.arg1 = ExpectedArg1
	testObject.arg2 = expectedArg2
	command, arg1, arg2, err := testObject.GetMessage()

	assert.String("command", ExpectedCommand, command)
	assert.String("arg1", ExpectedArg1, arg1)
	assert.String("arg2", expectedArg2, arg2)
	assert.Error(nil, err)
}
