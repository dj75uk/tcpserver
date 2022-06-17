package kvstore

import (
	"errors"
)

const kvCommandUpsert string = "UPSERT"
const kvCommandGet string = "GET"
const kvCommandDelete string = "DELETE"
const kvCommandList string = "LIST"

type KvStore struct {
	items    map[string]string
	requests chan kvStoreRequest
}

var ErrKeyNotFound error = errors.New("key not found")

type kvStoreRequest struct {
	Command string
	Key     string
	Value   string
	Results chan kvStoreResponse
}

type kvStoreResponse struct {
	Value  string
	Values []string
	Error  error
}

func NewKvStore() *KvStore {
	return &KvStore{
		items:    nil,
		requests: nil,
	}
}

func (store *KvStore) Open() {
	if store.items == nil {
		store.items = make(map[string]string)
		store.requests = make(chan kvStoreRequest)
		go handleRequests(store)
	}
}

func (store *KvStore) Close() {
	if store.items != nil {
		store.items = nil
		close(store.requests)
		store.requests = nil
	}
}

func (store *KvStore) query(request kvStoreRequest) {
	store.requests <- request
}

func (store *KvStore) Get(key string) (string, error) {
	request := kvStoreRequest{
		Command: kvCommandGet,
		Key:     key,
		Value:   "",
		Results: make(chan kvStoreResponse),
	}
	store.query(request)
	response := <-request.Results
	defer close(request.Results)
	return response.Value, response.Error
}

func (store *KvStore) Upsert(key string, value string) (string, error) {
	request := kvStoreRequest{
		Command: kvCommandUpsert,
		Key:     key,
		Value:   value,
		Results: make(chan kvStoreResponse),
	}
	store.query(request)
	response := <-request.Results
	defer close(request.Results)
	return response.Value, response.Error
}

func (store *KvStore) Delete(key string) (string, error) {
	request := kvStoreRequest{
		Command: kvCommandDelete,
		Key:     key,
		Value:   "",
		Results: make(chan kvStoreResponse),
	}
	store.query(request)
	response := <-request.Results
	defer close(request.Results)
	return response.Value, response.Error
}

func (store *KvStore) ListKeys() []string {
	request := kvStoreRequest{
		Command: kvCommandList,
		Key:     "",
		Value:   "",
		Results: make(chan kvStoreResponse),
	}
	store.query(request)
	response := <-request.Results
	defer close(request.Results)
	return response.Values
}

func handleRequests(store *KvStore) {
	for {
		request, isOpen := <-store.requests
		if !isOpen {
			return
		}
		response := kvStoreResponse{
			Value:  "",
			Values: nil,
			Error:  nil,
		}
		switch request.Command {
		case kvCommandUpsert:
			previous, exists := store.items[request.Key]
			if !exists {
				store.items[request.Key] = request.Value
			} else {
				if previous != request.Value {
					store.items[request.Key] = request.Value
				}
			}
		case kvCommandGet:
			value, exists := store.items[request.Key]
			if !exists {
				response.Error = ErrKeyNotFound
			} else {
				response.Value = value
			}
		case kvCommandDelete:
			delete(store.items, request.Key)
		case kvCommandList:
			response.Values = make([]string, 0)
			for k := range store.items {
				response.Values = append(response.Values, k)
			}
		}
		request.Results <- response
	}
}
