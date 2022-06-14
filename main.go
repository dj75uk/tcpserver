package main

import (
	"fmt"
	"kvsapp/kvserver"
	"kvsapp/kvstore"
	"os"
)

func main() {

	store := kvstore.NewKvStore()
	store.Open()
	defer store.Close()

	server, err := kvserver.NewKvServer(8000, store)
	if err != nil {
		os.Exit(-1)
	}

	err = server.Open()
	if err != nil {
		os.Exit(-2)
	}
	defer server.Close()

	foo := make(chan int)
	bar := <-foo
	fmt.Println(bar)
}
