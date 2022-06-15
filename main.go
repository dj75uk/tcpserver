package main

import (
	"kvsapp/kvserver"
	"kvsapp/kvstore"
	"os"
)

func main() {

	// create a new store...
	store := kvstore.NewKvStore()
	store.Open()
	defer store.Close()

	// create a new server...
	server, err := kvserver.NewKvServer(8000, store)
	if err != nil {
		os.Exit(-1)
	}

	// start the server...
	err = server.Open()
	if err != nil {
		os.Exit(-2)
	}
	defer server.Close()

	// wait for ctrl-c without killing the cpu...
	<-make(chan int)
}
