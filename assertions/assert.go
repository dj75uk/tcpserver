package assertions

import (
	"errors"
	"testing"
)

type Assert struct {
	t *testing.T
}

func NewAssert(t *testing.T) Assert {
	return Assert{t: t}
}

func (assert Assert) TestError(testName string, expected error, actual error) {
	if expected == nil && actual != nil {
		assert.t.Errorf("test: %s, param: err, expected: nil, actual: %T (%s)", testName, actual, actual.Error())
	}
	if expected != nil && actual == nil {
		assert.t.Errorf("test: %s, param: err, expected: %T (%s), actual: nil", testName, expected, expected.Error())
	}
	if expected != nil && actual != nil && !errors.Is(expected, actual) {
		assert.t.Errorf("test: %s, param: err, expected: %T (%s), actual: %T (%s)", testName, expected, expected.Error(), actual, actual.Error())
	}
}

func (assert Assert) Error(expected error, actual error) {
	if expected == nil && actual != nil {
		assert.t.Errorf("param: err, expected: nil, actual: %T (%s)", actual, actual.Error())
	}
	if expected != nil && actual == nil {
		assert.t.Errorf("param: err, expected: %T (%s), actual: nil", expected, expected.Error())
	}
	if expected != nil && actual != nil && !errors.Is(expected, actual) {
		assert.t.Errorf("param: err, expected: %T (%s), actual: %T (%s)", expected, expected.Error(), actual, actual.Error())
	}
}

func (assert Assert) TestBoolean(testName string, param string, expected bool, actual bool) {
	if expected != actual {
		assert.t.Errorf("test: %s, param: %s, expected: %v, actual: %v", testName, param, expected, actual)
	}
}
func (assert Assert) Boolean(param string, expected bool, actual bool) {
	if expected != actual {
		assert.t.Errorf("param: %s, expected: %v, actual: %v", param, expected, actual)
	}
}
func (assert Assert) True(param string, actual bool) {
	assert.Boolean(param, true, actual)
}
func (assert Assert) False(param string, actual bool) {
	assert.Boolean(param, false, actual)
}

func (assert Assert) TestString(testName string, param string, expected string, actual string) {
	if expected != actual {
		assert.t.Errorf("test: %s, param: %s, expected: %s, actual: %s", testName, param, expected, actual)
	}
}
func (assert Assert) String(param string, expected string, actual string) {
	if expected != actual {
		assert.t.Errorf("param: %s, expected: %s, actual: %s", param, expected, actual)
	}
}
