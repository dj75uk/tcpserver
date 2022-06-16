package kvserver

import (
	"context"
	"fmt"
	"kvsapp/parsing"
	"net"
	"os"
	"syscall"
	"time"
)

const udpNetwork string = "udp4"

func getServerHostKey() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("[%s:%d:%d]", hostname, os.Getppid(), os.Getpid())
}

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
	listenerConfiguration := getUdpListenerConfiguration()

	//kvs.udpListeningAddress = "192.168.0.106:9000"
	kvs.udpListeningAddress = ":9000"

	udpConnection, err := listenerConfiguration.ListenPacket(context.Background(), "udp", kvs.udpListeningAddress)
	if err != nil {
		fmt.Printf("server-udp: unable to start listening, err: %s\n", err.Error())
		return
	}
	fmt.Println("server-udp: listening for broadcasts")

	defer udpConnection.Close()

	buffer := make([]byte, 2000)
	for {
		if readCount, _, err := udpConnection.ReadFrom(buffer); err == nil && readCount > 0 {
			if parser, err := parsing.NewParser(getStandardGrammar()); err == nil {
				for _, v := range buffer {
					if found, err := parser.Process(string(v)); err == nil && found {
						if cmd, arg1, arg2, err := parser.GetMessage(); err == nil {
							if arg1 != hostKey {
								_ = kvs.handleMessage(nil, &commandMessage{Command: cmd, Key: arg1, Value: arg2})
							}
						}
						break
					}
				}
			}
		}
	}
}

func (kvs *KvServer) handleUdpBroadcast(hostKey string) {
	kvs.udpBroadcastAddress = "192.168.0.255:9000"

	broadcastAddress, _ := net.ResolveUDPAddr(udpNetwork, kvs.udpBroadcastAddress)
	for {
		time.Sleep(3 * time.Second)
		if udpConnection, err := net.DialUDP(udpNetwork, nil, broadcastAddress); err == nil {
			if txBuffer, err := parsing.CreateData("hst", getServerHostKey(), fmt.Sprintf("tcp=%d", kvs.tcpport)); err == nil {
				if count, err := udpConnection.Write(txBuffer); err == nil && count > 0 {
					fmt.Println("server-udp: broadcast")
				}
			}
			udpConnection.Close()
		}
	}
}
