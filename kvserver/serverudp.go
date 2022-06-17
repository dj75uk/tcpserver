package kvserver

import (
	"context"
	"fmt"
	"kvsapp/parsing"
	"net"
	"syscall"
	"time"
)

const udpNetwork string = "udp"

// creates a listener configuration suitable for broadcast
func getUdpListenerConfiguration() net.ListenConfig {
	return net.ListenConfig{
		Control: func(network string, address string, c syscall.RawConn) error {
			return c.Control(func(fd uintptr) {
				handle := syscall.Handle(fd)
				syscall.SetsockoptInt(handle, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				syscall.SetsockoptInt(handle, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1)
			})
		},
	}
}

func (kvs *KvServer) handleUdpListener(hostKey string) {

	//kvs.udpListeningAddress = "192.168.0.106:9000"
	kvs.udpListeningAddress = fmt.Sprintf("0.0.0.0:%d", kvs.udpport)
	//kvs.udpListeningAddress = ":9000"

	listenerConfiguration := getUdpListenerConfiguration()
	udpConnection, err := listenerConfiguration.ListenPacket(context.Background(), udpNetwork, kvs.udpListeningAddress)
	if err != nil {
		fmt.Printf("cluster: unable to start udp listening, err: %s\n", err.Error())
		return
	}
	fmt.Println("cluster: listening for broadcasts")

	defer udpConnection.Close()

	buffer := make([]byte, 2000)
	for {
		// wait for data...
		//udpConnection.SetReadDeadline(time.Now().Add(1 * time.Second))
		readCount, _, err := udpConnection.ReadFrom(buffer)
		if readCount <= 0 || err != nil {
			continue
		}

		// create a new parser for the message...
		msg := string(buffer[0:readCount])
		fmt.Println("cluster: incoming broadcast message")
		parser, err := parsing.NewParser(getStandardGrammar())
		if err != nil {
			continue
		}

		// attempt to parse the message...
		found := false
		for i := 0; i < readCount; i++ {
			found, err = parser.Process(string(msg[i]))
			if found && err == nil {
				break
			}
			if err != nil {
				found = false
				break
			}
		}
		if !found {
			continue
		}

		// obtain the message object...
		cmd, arg1, arg2, err := parser.GetMessage()
		//fmt.Printf("cluster: GetMessage() cmd: %s, arg1: %s, arg2: %s, err: %v\n", cmd, arg1, arg2, err)
		if err != nil {
			continue
		}

		// remove messages from ourselves...
		if arg1 == hostKey {
			fmt.Printf("cluster: skipping message because it's from us!\n")
			continue
		}

		// process the message via the standard message handler...
		_ = kvs.handleMessage(nil, &commandMessage{Command: cmd, Key: arg1, Value: arg2})

	}
}

func (kvs *KvServer) handleUdpBroadcast(hostKey string, tcpAddress string) {
	kvs.udpBroadcastAddress = fmt.Sprintf("192.168.0.255:%d", kvs.udpport)
	//kvs.udpBroadcastAddress = "0.0.0.255:9000"

	broadcastAddress, _ := net.ResolveUDPAddr(udpNetwork, kvs.udpBroadcastAddress)
	for {
		// wait an arbitrary time to avoid network spamming...
		time.Sleep(3 * time.Second)

		// create a udp broadcast connection...
		udpConnection, err := net.DialUDP(udpNetwork, nil, broadcastAddress)
		if err != nil {
			continue
		}

		// create a message to broadcast...
		txBuffer, err := parsing.CreateData("hst", getServerHostKey(), tcpAddress)
		if err != nil {
			continue
		}

		// broadcast the message...
		if count, err := udpConnection.Write(txBuffer); err == nil && count > 0 {
			fmt.Println("cluster: broadcasting identity")
		}

		// close the broadcast connection...
		udpConnection.Close()
	}
}
