package kvstore_test

import (
	"kvsapp/assertions"
	"kvsapp/kvstore"
	"testing"
)

func createTestObject() *kvstore.KvStore {
	return kvstore.NewKvStore()
}

func TestNewStoreReturnsObject(t *testing.T) {
	t.Parallel()
	testObject := createTestObject()
	if testObject == nil {
		t.Errorf("expected: obj, actual: nil")
	}
}

func TestGetReturnsErrorOnUnknownKey(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	store := createTestObject()
	store.Open()
	defer store.Close()

	expectedKey := "TestGetReturnsErrorOnUnknownKey"
	actualValue, err := store.Get(expectedKey)
	assert.String("value", "", actualValue)
	assert.Error(kvstore.ErrKeyNotFound, err)
}

func TestUpsertInsertsNewItem(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	store := createTestObject()
	store.Open()
	defer store.Close()

	expectedKey := "TestUpsertInsertsNewItem"
	expectedValue := "TestUpsertInsertsNewItemValue"

	store.Upsert(expectedKey, expectedValue)

	actualValue, err := store.Get(expectedKey)
	assert.String("value", expectedValue, actualValue)
	assert.Error(nil, err)
}

func TestUpsertUpdatesExistingItem(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	store := createTestObject()
	store.Open()
	defer store.Close()

	expectedKey := "TestUpsertUpdatesExistingItem"
	expectedValue1 := "TestUpsertUpdatesExistingItemValue1"
	expectedValue2 := "TestUpsertUpdatesExistingItemValue2"

	store.Upsert(expectedKey, expectedValue1)
	store.Upsert(expectedKey, expectedValue2)

	actualValue, err := store.Get(expectedKey)
	assert.String("value", expectedValue2, actualValue)
	assert.Error(nil, err)
}

func TestDeleteRemovesKey(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	store := createTestObject()
	store.Open()
	defer store.Close()

	expectedKey := "TestDeleteRemovesKey"
	expectedValue := "TestDeleteRemovesKeyValue1"

	if _, err := store.Upsert(expectedKey, expectedValue); err != nil {
		t.Fatalf("test setup failure (upsert)")
	}
	if _, err := store.Delete(expectedKey); err != nil {
		t.Fatalf("test setup failure (delete)")
	}
	actualValue, err := store.Get(expectedKey)
	assert.String("value", "", actualValue)
	assert.Error(kvstore.ErrKeyNotFound, err)
}

func TestDeleteReturnsNoErrorOnUnknownKey(t *testing.T) {
	t.Parallel()
	assert := assertions.NewAssert(t)
	store := createTestObject()
	store.Open()
	defer store.Close()

	expectedKey := "TestGetReturnsErrorOnUnknownKey"
	actualValue, err := store.Delete(expectedKey)
	assert.String("value", "", actualValue)
	assert.Error(nil, err)
}
