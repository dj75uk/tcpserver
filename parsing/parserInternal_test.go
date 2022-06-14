package parsing

import (
	"errors"
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
	if testObject.state != expectedState {
		t.Errorf("param: %s, expected: %d, actual: %d", "state", expectedState, testObject.state)
	}
	if testObject.command != expectedCommand {
		t.Errorf("param: %s, expected: %s, actual: %s", "command", expectedCommand, testObject.command)
	}
	if testObject.argsExpected != expectedArgsExpected {
		t.Errorf("param: %s, expected: %d, actual: %d", "argsExpected", expectedArgsExpected, testObject.argsExpected)
	}
}

func assertArg1(t *testing.T, testObject *Parser, expectedArg1LengthLength int, expectedArg1LengthBuilder string, expectedArg1Length int, expectedArg1 string) {
	if testObject.arg1LengthLength != expectedArg1LengthLength {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg1LengthLength", expectedArg1LengthLength, testObject.arg1LengthLength)
	}
	if testObject.arg1LengthBuilder != expectedArg1LengthBuilder {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg1LengthBuilder", expectedArg1LengthBuilder, testObject.arg1LengthBuilder)
	}
	if testObject.arg1Length != expectedArg1Length {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg1Length", expectedArg1Length, testObject.arg1Length)
	}
	if testObject.arg1 != expectedArg1 {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg1", expectedArg1, testObject.arg1)
	}
}

func assertArg2(t *testing.T, testObject *Parser, expectedArg2LengthLength int, expectedArg2LengthBuilder string, expectedArg2Length int, expectedArg2 string) {
	if testObject.arg2LengthLength != expectedArg2LengthLength {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg2LengthLength", expectedArg2LengthLength, testObject.arg2LengthLength)
	}
	if testObject.arg2LengthBuilder != expectedArg2LengthBuilder {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg2LengthBuilder", expectedArg2LengthBuilder, testObject.arg2LengthBuilder)
	}
	if testObject.arg2Length != expectedArg2Length {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg2Length", expectedArg2Length, testObject.arg2Length)
	}
	if testObject.arg2 != expectedArg2 {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg2", expectedArg2, testObject.arg2)
	}
}

func assertResetState(t *testing.T, testObject *Parser, testGrammar map[string]uint16) {
	assertState(t, testObject, 0, "", 0)
	assertArg1(t, testObject, 0, "", 0, "")
	assertArg2(t, testObject, 0, "", 0, "")
	assertGrammar(t, testObject, testGrammar)
}

func assertError(t *testing.T, expected error, actual error) {
	if expected == nil && actual != nil {
		t.Errorf("param: err, expected: nil, actual: %T (%s)", actual, actual.Error())
	}
	if expected != nil && actual == nil {
		t.Errorf("param: err, expected: %T (%s), actual: nil", expected, expected.Error())
	}
	if expected != nil && actual != nil && !errors.Is(expected, actual) {
		t.Errorf("param: err, expected: %T (%s), actual: %T (%s)", expected, expected.Error(), actual, actual.Error())
	}
}

func assertBoolean(t *testing.T, param string, expected bool, actual bool) {
	if expected != actual {
		t.Errorf("param: %s, expected: %v, actual: %v", param, expected, actual)
	}
}
func assertTrue(t *testing.T, param string, actual bool) {
	assertBoolean(t, param, true, actual)
}
func assertFalse(t *testing.T, param string, actual bool) {
	assertBoolean(t, param, false, actual)
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

func TestProcess0001(t *testing.T) {
	t.Parallel()
	testGrammar := map[string]uint16{
		"ddd": 0,
		"eee": 1,
		"fff": 2,
	}
	testObject, _ := NewParser(testGrammar)
	assertResetState(t, testObject, testGrammar)

	found := false
	err := error(nil)

	found, err = testObject.Process("x")
	assertFalse(t, "found", found)
	assertError(t, nil, err)
	assertState(t, testObject, 0, "x", 0)

	found, err = testObject.Process("y")
	assertFalse(t, "found", found)
	assertError(t, nil, err)
	assertState(t, testObject, 0, "xy", 0)

	found, err = testObject.Process("z")
	assertFalse(t, "found", found)
	assertError(t, ErrParserUnknownCommand, err)
	assertState(t, testObject, 0, "", 0)
}
