package main

import (
	"flag"
	"fmt"
	"kvsapp/kvserver"
	"kvsapp/kvstore"
	"os"
)

const DefaultPortNumber int = 8000

func main() {
	var port = DefaultPortNumber
	flag.IntVar(&port, "port", DefaultPortNumber, "port number to listen on")
	flag.Parse()

	// create a new store...
	store := kvstore.NewKvStore()
	store.Open()
	defer store.Close()

	// create a new server...
	server, err := kvserver.NewKvServer(port, store)
	if err != nil {
		os.Exit(-1)
	}
	fmt.Printf("server listening on port %d\n", port)

	// start the server...
	err = server.Open()
	if err != nil {
		os.Exit(-2)
	}
	defer server.Close()

	// wait for ctrl-c without killing the cpu...
	<-make(chan int)
}
