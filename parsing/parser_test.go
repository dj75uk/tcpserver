package parsing_test

import (
	"kvsapp/parsing"
	"testing"
)

type testItem struct {
	bytes  []byte
	params uint16
	expCmd string
	expKey string
	expVal string
	expMsg bool
	expPro int
	expErr bool
}

func getTestData() map[string]testItem {
	result := map[string]testItem{
		"test00": {params: 0, bytes: []byte(""), expCmd: "", expKey: "", expVal: "", expMsg: false, expPro: 0, expErr: false},
		"test01": {params: 0, bytes: []byte("cmd"), expCmd: "cmd", expKey: "", expVal: "", expMsg: true, expPro: 3, expErr: false},
		"test02": {params: 0, bytes: []byte("cmdIGNORE"), expCmd: "cmd", expKey: "", expVal: "", expMsg: true, expPro: 3, expErr: false},
		"test03": {params: 1, bytes: []byte("cmd13key"), expCmd: "cmd", expKey: "key", expVal: "", expMsg: true, expPro: 8, expErr: false},
		"test04": {params: 1, bytes: []byte("cmd13keyIGNORE"), expCmd: "cmd", expKey: "key", expVal: "", expMsg: true, expPro: 8, expErr: false},
		"test05": {params: 2, bytes: []byte("cmd13key15value"), expCmd: "cmd", expKey: "key", expVal: "value", expMsg: true, expPro: 15, expErr: false},
		"test06": {params: 2, bytes: []byte("cmd13key15valueIGNORE"), expCmd: "cmd", expKey: "key", expVal: "value", expMsg: true, expPro: 15, expErr: false},
		"test07": {params: 1, bytes: []byte("cmd03key"), expCmd: "cmd", expKey: "", expVal: "", expMsg: false, expPro: 4, expErr: true},
		"test08": {params: 2, bytes: []byte("cmd13key05value"), expCmd: "", expKey: "", expVal: "", expMsg: false, expPro: 9, expErr: true},
	}

	return result
}

func TestParse(t *testing.T) {
	for testName, testData := range getTestData() {

		testObject := createTestObject()
		actMsg, actCmp, actErr := testObject.Parse(testData.bytes, testData.params)

		if actMsg == nil && testData.expMsg {
			t.Errorf("test: %s, param: message, expected: object, actual: nil", testName)
		}
		if actMsg != nil && !testData.expMsg {
			t.Errorf("test: %s, param: message, expected: nil, actual: object", testName)
		}
		if actCmp != testData.expPro {
			t.Errorf("test: %s, param: complete, expected: %v, actual: %v", testName, testData.expPro, actCmp)
		}
		if actErr == nil && testData.expErr {
			t.Errorf("test: %s, param: error, expected: object, actual: nil", testName)
		}
		if actErr != nil && !testData.expErr {
			t.Errorf("test: %s, param: error, expected: nil, actual: object", testName)
		}

		if actMsg != nil && testData.expMsg {

			if actMsg.Command != testData.expCmd {
				t.Errorf("test: %s, param: message.command, expected: %s, actual: %s", testName, testData.expCmd, actMsg.Command)
			}
			if actMsg.Key != testData.expKey {
				t.Errorf("test: %s, param: message.key, expected: %s, actual: %s", testName, testData.expKey, actMsg.Key)
			}
			if actMsg.Value != testData.expVal {
				t.Errorf("test: %s, param: message.value, expected: %s, actual: %s", testName, testData.expVal, actMsg.Value)
			}

		}
	}
}

func createTestObject() *parsing.Parser {
	return parsing.NewParser()
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

func TestNewParser(t *testing.T) {
	p := parsing.NewParser()
	if p == nil {
		t.Error("expected: object, actual: nil")
	}
}

func TestCreateDataOnEmptyCommandErrorIsReturned(t *testing.T) {
	p := createTestObject()
	_, err := p.CreateData("", "", "")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataOnTooSmallCommandErrorIsReturned(t *testing.T) {
	p := createTestObject()
	_, err := p.CreateData("a", "", "")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataOnTooLargeCommandErrorIsReturned(t *testing.T) {
	p := createTestObject()
	_, err := p.CreateData("aaaa", "", "")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataOnCorrectCommandNoErrorIsReturned(t *testing.T) {
	p := createTestObject()
	_, err := p.CreateData("aaa", "", "")
	if err != nil {
		t.Errorf("expected: nil, actual: error (%s)", err.Error())
	}
}

func TestCreateDataOnValueWithNoKeyErrorIsReturned(t *testing.T) {
	p := createTestObject()
	_, err := p.CreateData("aaa", "", "vvvvvvv")
	if err == nil {
		t.Error("expected: error, actual: nil")
	}
}

func TestCreateDataWithNoKeyOrValue(t *testing.T) {
	p := createTestObject()
	expected := []byte("CMD")
	actual, _ := p.CreateData("CMD", "", "")
	if !compareSlices(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestCreateDataWithKey(t *testing.T) {
	p := createTestObject()
	expected := []byte("CMD17KEYNAME")
	actual, _ := p.CreateData("CMD", "KEYNAME", "")
	if !compareSlices(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestCreateDataWithKeyAndValue(t *testing.T) {
	p := createTestObject()
	expected := []byte("CMD17KEYNAME220SOME ARBITRARY VALUE")
	actual, _ := p.CreateData("CMD", "KEYNAME", "SOME ARBITRARY VALUE")
	if !compareSlices(actual, expected) {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func TestParseOnEmptyStringIsCorrect(t *testing.T) {
	p := createTestObject()
	msg, complete, err := p.Parse([]byte(""), 0)
	if msg != nil {
		t.Error("expected: nil, actual: object (msg)")
	}
	if complete != 0 {
		t.Errorf("expected: %d, actual: %d (processed)", 0, complete)
	}
	if err != nil {
		t.Error("expected: nil, actual: object (error)")
	}
}

func TestParseP0MessagesResultsInCorrectCommand(t *testing.T) {
	c := "abc"
	p := createTestObject()
	msg, _, _ := p.Parse([]byte(c), 0)
	if msg == nil {
		t.Error("expected: nil, actual: object (msg)")
		return
	}
	if msg.Command != c {
		t.Errorf("expected: %s, actual: %s", c, msg.Command)
	}
}

func TestParseP0MessagesResultsInCorrectKey(t *testing.T) {
	c := "abc"
	p := createTestObject()
	msg, _, _ := p.Parse([]byte(c), 0)
	if msg == nil {
		t.Error("expected: nil, actual: object (msg)")
		return
	}
	if msg.Key != "" {
		t.Errorf("expected: %s, actual: %s", "", msg.Key)
	}
}

func TestParseP0MessagesResultsInCorrectValue(t *testing.T) {
	c := "abc"
	p := createTestObject()
	msg, _, _ := p.Parse([]byte(c), 0)
	if msg == nil {
		t.Error("expected: nil, actual: object (msg)")
		return
	}
	if msg.Value != "" {
		t.Errorf("expected: %s, actual: %s", "", msg.Value)
	}
}

func TestParseP1MessagesResultsInCorrectCommand(t *testing.T) {
	c := "abc"
	k := "keyname"
	v := ""
	p := createTestObject()
	data, _ := p.CreateData(c, k, v)
	msg, _, _ := p.Parse(data, 1)
	if msg == nil {
		t.Error("expected: nil, actual: object (msg)")
		return
	}
	if msg.Command != c {
		t.Errorf("expected: %s, actual: %s", k, msg.Command)
	}
}

func TestParseP1MessagesResultsInCorrectKey(t *testing.T) {
	c := "abc"
	k := "keyname"
	v := ""
	p := createTestObject()
	data, _ := p.CreateData(c, k, v)
	msg, _, _ := p.Parse(data, 1)
	if msg == nil {
		t.Error("expected: nil, actual: object (msg)")
		return
	}
	if msg.Key != k {
		t.Errorf("expected: %s, actual: %s", k, msg.Key)
	}
}

func TestParseP1MessagesResultsInCorrectValue(t *testing.T) {
	c := "abc"
	k := "keyname"
	v := ""
	p := createTestObject()
	data, _ := p.CreateData(c, k, v)
	msg, _, _ := p.Parse(data, 1)
	if msg == nil {
		t.Error("expected: nil, actual: object (msg)")
		return
	}
	if msg.Value != v {
		t.Errorf("expected: %s, actual: %s", k, msg.Value)
	}
}
