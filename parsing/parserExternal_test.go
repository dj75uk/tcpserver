package parsing_test

import (
	"errors"
	"kvsapp/assertions"
	"kvsapp/parsing"
	"sync"
	"testing"
)

type sampleData struct {
	enabled              bool
	bytes                []byte
	expectedCommand      string
	expectedArg1         string
	expectedArg2         string
	expectedFound        bool
	expectedProcessError error
}

func createTestObject() *parsing.Parser {
	result, _ := parsing.NewParser(map[string]uint16{
		"cm0": 0,
		"cm1": 1,
		"cm2": 2,
	})
	return result
}

func getSampleData() map[string]sampleData {
	result := map[string]sampleData{
		"0 bytes": {
			enabled:              true,
			bytes:                []byte(""),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: nil,
		},
		"1 byte": {
			enabled:              true,
			bytes:                []byte("c"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: nil,
		},
		"2 bytes": {
			enabled:              true,
			bytes:                []byte("cm"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: nil,
		},
		"0-args": {
			enabled:              true,
			bytes:                []byte("cm0"),
			expectedCommand:      "cm0",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        true,
			expectedProcessError: nil,
		},
		"0-args-skips": {
			enabled:              true,
			bytes:                []byte("cm0xxx"),
			expectedCommand:      "cm0",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        true,
			expectedProcessError: nil,
		},
		"1-args": {
			enabled:              true,
			bytes:                []byte("cm114arg1"),
			expectedCommand:      "cm1",
			expectedArg1:         "arg1",
			expectedArg2:         "",
			expectedFound:        true,
			expectedProcessError: nil,
		},
		"1-args-skips": {
			enabled:              true,
			bytes:                []byte("cm114arg1xxx"),
			expectedCommand:      "cm1",
			expectedArg1:         "arg1",
			expectedArg2:         "",
			expectedFound:        true,
			expectedProcessError: nil,
		},
		"2-args": {
			enabled:              true,
			bytes:                []byte("cm214arg114arg2"),
			expectedCommand:      "cm2",
			expectedArg1:         "arg1",
			expectedArg2:         "arg2",
			expectedFound:        true,
			expectedProcessError: nil,
		},
		"2-args-skips": {
			enabled:              true,
			bytes:                []byte("cm214arg114arg2xxx"),
			expectedCommand:      "cm2",
			expectedArg1:         "arg1",
			expectedArg2:         "arg2",
			expectedFound:        true,
			expectedProcessError: nil,
		},
		"bad format (1)": {
			enabled:              true,
			bytes:                []byte("cm10"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (2)": {
			enabled:              true,
			bytes:                []byte("cm1x"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (3)": {
			enabled:              true,
			bytes:                []byte("cm110"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (4)": {
			enabled:              true,
			bytes:                []byte("cm11x"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (5)": {
			enabled:              true,
			bytes:                []byte("cm211k0"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (6)": {
			enabled:              true,
			bytes:                []byte("cm211kx"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (7)": {
			enabled:              true,
			bytes:                []byte("cm211k10"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"bad format (8)": {
			enabled:              true,
			bytes:                []byte("cm211k1x"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserBadFormat,
		},
		"unknown command": {
			enabled:              true,
			bytes:                []byte("xxx"),
			expectedCommand:      "",
			expectedArg1:         "",
			expectedArg2:         "",
			expectedFound:        false,
			expectedProcessError: parsing.ErrParserUnknownCommand,
		},
	}

	return result
}

func TestSampleData(t *testing.T) {
	t.Parallel()
	wait := sync.WaitGroup{}
	for testName, testData := range getSampleData() {
		if !testData.enabled {
			continue
		}

		wait.Add(1)

		go func(t *testing.T, testName string, testData sampleData) {
			defer wait.Done()
			testObject := createTestObject()

			processFound := false
			processError := error(nil)

			for _, testByte := range testData.bytes {

				processFound, processError = testObject.Process(string(testByte))
				if processFound {
					break
				}
				if processError != nil {
					break
				}
			}

			if processFound != testData.expectedFound {
				t.Errorf("test: %s, param: processFound, expected: %v, actual: %v", testName, testData.expectedFound, processFound)
			}

			if (processError != nil) && (testData.expectedProcessError != nil) {
				if !errors.Is(processError, testData.expectedProcessError) {
					t.Errorf("test1: %s, param: processError, expected: %T (%v), actual: %T (%v)", testName, testData.expectedProcessError, testData.expectedProcessError, processError, processError)
				}
			}
			if (processError == nil) && (testData.expectedProcessError != nil) {
				t.Errorf("test2: %s, param: processError, expected: %T (%v), actual: nil", testName, testData.expectedProcessError, testData.expectedProcessError)
			}
			if (processError != nil) && (testData.expectedProcessError == nil) {
				t.Errorf("test2: %s, param: processError, expected: nil, actual: %T (%v)", testName, processError, processError)
			}
			if !processFound || processError != nil {
				return
			}

			getMessageCommand, getMessageArg1, getMessageArg2, getMessageError := testObject.GetMessage()

			if getMessageCommand != testData.expectedCommand {
				t.Errorf("test: %s, param: command, expected: %s, actual: %s", testName, testData.expectedCommand, getMessageCommand)
			}
			if getMessageArg1 != testData.expectedArg1 {
				t.Errorf("test: %s, param: arg1, expected: %s, actual: %s", testName, testData.expectedArg1, getMessageArg1)
			}
			if getMessageArg2 != testData.expectedArg2 {
				t.Errorf("test: %s, param: arg2, expected: %s, actual: %s", testName, testData.expectedArg2, getMessageArg2)
			}
			if getMessageError != nil {
				t.Errorf("test: %s, param: error, expected: %s, actual: %s", testName, testData.expectedCommand, getMessageCommand)
			}
		}(t, testName, testData)
	}
	wait.Wait()
}

func compareSlices(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGetMessageReturnsErrorWhenNoMessageReady(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	testObject, _ := parsing.NewParser(map[string]uint16{"cmd": 0})
	testObject.Process("a")
	getMessageCommand, getMessageArg1, getMessageArg2, getMessageError := testObject.GetMessage()
	assert.String("command", "", getMessageCommand)
	assert.String("arg1", "", getMessageArg1)
	assert.String("arg2", "", getMessageArg2)
	assert.Error(parsing.ErrParserNoMessage, getMessageError)
}

func TestNewParserReturnsErrorOnNilArgument(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	testObject, err := parsing.NewParser(nil)
	if testObject != nil {
		t.Errorf("param: %s, expected: %s, actual: %T", "testObject", "nil", testObject)
	}
	assert.Error(parsing.ErrParserInvalidArgument, err)
}

func TestNewParserReturnsObject(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	testObject, err := parsing.NewParser(map[string]uint16{"cmd": 0})
	if testObject == nil {
		t.Errorf("param: %s, expected: %s, actual: %s", "testObject", "obj", "nil")
	}
	assert.Error(nil, err)
}

func TestCreateDataOnEmptyCommandErrorIsReturned(t *testing.T) {
	t.Parallel()
	_, err := parsing.CreateData("", "", "")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataOnTooSmallCommandErrorIsReturned(t *testing.T) {
	t.Parallel()
	_, err := parsing.CreateData("a", "", "")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataOnTooLargeCommandErrorIsReturned(t *testing.T) {
	t.Parallel()
	_, err := parsing.CreateData("aaaa", "", "")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataOnCorrectCommandNoErrorIsReturned(t *testing.T) {
	t.Parallel()
	_, err := parsing.CreateData("aaa", "", "")
	if err != nil {
		t.Errorf("expected: nil, actual: error (%s)", err.Error())
	}
}

func TestCreateDataOnValueWithNoKeyErrorIsReturned(t *testing.T) {
	t.Parallel()
	_, err := parsing.CreateData("aaa", "", "vvvvvvv")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataWithNoKeyOrValue(t *testing.T) {
	t.Parallel()
	expected := []byte("CMD")
	actual, _ := parsing.CreateData("CMD", "", "")
	if !compareSlices(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestCreateDataWithKey(t *testing.T) {
	t.Parallel()
	expected := []byte("CMD17KEYNAME")
	actual, _ := parsing.CreateData("CMD", "KEYNAME", "")
	if !compareSlices(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestCreateDataWithKeyAndValue(t *testing.T) {
	t.Parallel()
	expected := []byte("CMD17KEYNAME220SOME ARBITRARY VALUE")
	actual, _ := parsing.CreateData("CMD", "KEYNAME", "SOME ARBITRARY VALUE")
	if !compareSlices(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

// benchmark processing a fixed zero-argument message...
func BenchmarkProcessZeroArgumentMessage(b *testing.B) {
	testObject := createTestObject()
	for i := 0; i < b.N; i++ {
		_, _ = testObject.Process("c")
		_, _ = testObject.Process("m")
		_, _ = testObject.Process("0")
		_, _, _, _ = testObject.GetMessage()
	}
}

// benchmark processing a fixed one-argument message...
func BenchmarkProcessOneArgumentMessage(b *testing.B) {
	testObject := createTestObject()
	for i := 0; i < b.N; i++ {
		_, _ = testObject.Process("c")
		_, _ = testObject.Process("m")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("a")
		_, _, _, _ = testObject.GetMessage()
	}
}

// benchmark processing a fixed two-argument message...
func BenchmarkProcessTwoArgumentMessage(b *testing.B) {
	testObject := createTestObject()
	for i := 0; i < b.N; i++ {
		_, _ = testObject.Process("c")
		_, _ = testObject.Process("m")
		_, _ = testObject.Process("2")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("a")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("1")
		_, _ = testObject.Process("b")
		_, _, _, _ = testObject.GetMessage()
	}
}
