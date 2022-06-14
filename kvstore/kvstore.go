package kvstore

import (
	"errors"
	"fmt"
)

const KvCommandUpsert string = "UPSERT"
const KvCommandGet string = "GET"
const KvCommandDelete string = "DELETE"

type KvStore struct {
	items    map[string]string
	requests chan KvStoreRequest
}

type KvStoreRequest struct {
	Command string
	Key     string
	Value   string
	Results chan KvStoreResponse
}

type KvStoreResponse struct {
	Value string
	Error error
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
		store.requests = make(chan KvStoreRequest)
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

func (store *KvStore) Query(request KvStoreRequest) {
	store.requests <- request
}

func (store *KvStore) Get(key string) (string, error) {
	request := KvStoreRequest{
		Command: KvCommandGet,
		Key:     key,
		Value:   "",
		Results: make(chan KvStoreResponse),
	}
	store.Query(request)
	response := <-request.Results
	close(request.Results)
	return response.Value, response.Error
}

func (store *KvStore) Upsert(key string, value string) (string, error) {
	request := KvStoreRequest{
		Command: KvCommandUpsert,
		Key:     key,
		Value:   value,
		Results: make(chan KvStoreResponse),
	}
	store.Query(request)
	response := <-request.Results
	close(request.Results)
	return response.Value, response.Error
}

func (store *KvStore) Delete(key string) (string, error) {
	request := KvStoreRequest{
		Command: KvCommandDelete,
		Key:     key,
		Value:   "",
		Results: make(chan KvStoreResponse),
	}
	store.Query(request)
	response := <-request.Results
	close(request.Results)
	return response.Value, response.Error
}

func handleRequests(store *KvStore) {
	for {
		request, isOpen := <-store.requests
		if !isOpen {
			return
		}

		fmt.Printf("store:  REQ: {%s '%s' '%v'}\n", request.Command, request.Key, request.Value)

		response := KvStoreResponse{
			Value: "",
			Error: nil,
		}
		switch request.Command {
		case KvCommandUpsert:
			previous, exists := store.items[request.Key]
			if !exists {
				store.items[request.Key] = request.Value
			} else {
				if previous != request.Value {
					store.items[request.Key] = request.Value
				}
			}
		case KvCommandGet:
			value, exists := store.items[request.Key]
			if !exists {
				response.Error = errors.New("key not found")
			} else {
				response.Value = value
			}
		case KvCommandDelete:
			delete(store.items, request.Key)
		}
		fmt.Printf("store:  RES: {'%s' '%v'}\n", response.Value, response.Error)
		request.Results <- response
	}
}
