package main

import (
	"flag"
	"fmt"
	"kvsapp/kvserver"
	"kvsapp/kvstore"
	"os"
)

const DefaultTcpPortNumber int = 8000
const DefaultUdpPortNumber int = 9000

func main() {
	var tcpport = DefaultTcpPortNumber
	var udpport = DefaultUdpPortNumber
	flag.IntVar(&tcpport, "port", DefaultTcpPortNumber, "tcp port number to listen on")
	flag.IntVar(&udpport, "udpport", DefaultUdpPortNumber, "udp port number to listen on")
	flag.Parse()

	// create a new store...
	store := kvstore.NewKvStore()
	store.Open()
	defer store.Close()

	// create a new server...
	server, err := kvserver.NewKvServer(tcpport, udpport, store)
	if err != nil {
		fmt.Printf("server: error '%s'\n", err.Error())
		os.Exit(-1)
	}

	// start the server...
	err = server.Open()
	if err != nil {
		fmt.Printf("server: error '%s'\n", err.Error())
		os.Exit(-2)
	}
	defer server.Close()

	// wait for ctrl-c or server shutdown...
	server.WaitForShutdown()
}
