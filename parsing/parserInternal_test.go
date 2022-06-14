package parsing

import (
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

func assertResetState(t *testing.T, testObject *Parser, testGrammar map[string]uint16) {
	if testObject.state != 0 {
		t.Errorf("param: %s, expected: %d, actual: %d", "state", 0, testObject.state)
	}
	if testObject.command != "" {
		t.Errorf("param: %s, expected: %s, actual: %s", "command", "", testObject.command)
	}
	if testObject.argsExpected != 0 {
		t.Errorf("param: %s, expected: %d, actual: %d", "argsExpected", 0, testObject.argsExpected)
	}
	if testObject.arg1LengthLength != 0 {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg1LengthLength", 0, testObject.arg1LengthLength)
	}
	if testObject.arg1LengthBuilder != "" {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg1LengthBuilder", "", testObject.arg1LengthBuilder)
	}
	if testObject.arg1Length != 0 {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg1Length", 0, testObject.arg1Length)
	}
	if testObject.arg1 != "" {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg1", "", testObject.arg1)
	}
	if testObject.arg2LengthLength != 0 {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg2LengthLength", 0, testObject.arg2LengthLength)
	}
	if testObject.arg2LengthBuilder != "" {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg2LengthBuilder", "", testObject.arg2LengthBuilder)
	}
	if testObject.arg2Length != 0 {
		t.Errorf("param: %s, expected: %d, actual: %d", "arg2Length", 0, testObject.arg2Length)
	}
	if testObject.arg2 != "" {
		t.Errorf("param: %s, expected: %s, actual: %s", "arg2", "", testObject.arg2)
	}
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
