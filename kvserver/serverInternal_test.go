package kvserver

import (
	"bytes"
	"kvsapp/assertions"
	"kvsapp/kvstore"
	"kvsapp/parsing"
	"sync"
	"testing"
)

func createTestObject() *KvServer {
	store := kvstore.NewKvStore()
	store.Open()
	result, _ := NewKvServer(0, 0, store)
	return result
}

type handleReceivedByteTestData struct {
	bytes           string
	expectedCarryOn bool
	expectedError   error
	expectedWrite   string
}

func createHandleReceivedByteTestData() map[string]handleReceivedByteTestData {
	return map[string]handleReceivedByteTestData{
		"1 byte":           {bytes: "x", expectedCarryOn: true, expectedError: nil, expectedWrite: ""},
		"2 byte":           {bytes: "xy", expectedCarryOn: true, expectedError: nil, expectedWrite: ""},
		"unknown command":  {bytes: "xyz", expectedCarryOn: true, expectedError: nil, expectedWrite: "err"},
		"zero-arg command": {bytes: "bye", expectedCarryOn: false, expectedError: nil, expectedWrite: ""},
	}
}
func TestHandleReceivedByte(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	testObject := createTestObject()
	wait := sync.WaitGroup{}

	for testName, testData := range createHandleReceivedByteTestData() {
		wait.Add(1)
		go func(testName string, testData handleReceivedByteTestData) {
			defer wait.Done()
			parser, _ := parsing.NewParser(map[string]parsing.ParserGrammar{
				"bye": {ExpectedArguments: 0},
				"get": {ExpectedArguments: 1},
				"del": {ExpectedArguments: 1},
				"put": {ExpectedArguments: 2},
				"hed": {ExpectedArguments: 2},
			})
			buffer := &bytes.Buffer{}

			var carryOn bool
			var err error
			for _, b := range testData.bytes {
				carryOn, err = testObject.handleReceivedByte(buffer, parser, byte(b))
			}
			assert.TestError(testName, testData.expectedError, err)
			assert.TestBoolean(testName, "carryOn", testData.expectedCarryOn, carryOn)
			assert.TestString(testName, "written", testData.expectedWrite, buffer.String())
		}(testName, testData)
	}
	wait.Wait()
}
