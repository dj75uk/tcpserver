package kvserver

import (
	"errors"
	"fmt"
	"io"
	"kvsapp/kvstore"
	"kvsapp/parsing"
	"net"
	"os"
	"time"
)

const KvServerReadBufferSize int = 256

type KvServer struct {
	tcpport             int
	udpport             int
	udpBroadcastAddress string
	udpListeningAddress string
	store               *kvstore.KvStore
	servers             *kvstore.KvStore
	grammar             map[string]parsing.ParserGrammar
	shutdown            chan int
	handlers            map[string]func(kvs *KvServer, key string, value string) string
}

type commandMessage struct {
	Command string
	Key     string
	Value   string
}

func NewKvServer(tcpport int, udpport int, store *kvstore.KvStore) (*KvServer, error) {
	if store == nil {
		return nil, errors.New("parameter 'store' must not be nil")
	}
	return &KvServer{
		tcpport:  tcpport,
		udpport:  udpport,
		store:    store,
		servers:  kvstore.NewKvStore(),
		grammar:  getStandardGrammar(),
		shutdown: make(chan int),
		handlers: getHandlers(),
	}, nil
}

func (kvs *KvServer) Open() error {

	kvs.servers.Open()
	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", kvs.tcpport))
	if err != nil {
		return err
	}
	fmt.Printf("server-tcp: listening on port %d\n", kvs.tcpport)
	go kvs.handleTcpAcceptance(listener)

	tcpAddress := listener.Addr().String()

	go kvs.handleInternalChecking()
	go kvs.handleUdpListener(getServerHostKey())
	go kvs.handleUdpBroadcast(getServerHostKey(), tcpAddress)

	return nil
}

func getServerHostKey() string {
	hostname, _ := os.Hostname()
	return fmt.Sprintf("[%s:%d:%d]", hostname, os.Getppid(), os.Getpid())
}

func (kvs *KvServer) sendToAllOthers(command string, arg1 string, arg2 string) {
	for _, serverKey := range kvs.servers.ListKeys() {
		fault := true
		if serverAddress, err := kvs.servers.Get(serverKey); err == nil {
			if connection, err := net.Dial("tcp4", serverAddress); err == nil {
				data, _ := parsing.CreateData(command, arg1, arg2)
				connection.SetDeadline(time.Now().Add(800 * time.Millisecond))
				if written, err := connection.Write(data); written > 0 && err == nil {
					readBuffer := make([]byte, 16)
					if read, err := connection.Read(readBuffer); read > 0 && err == nil {
						fmt.Printf("cluster: sending '%s' command to '%s'\n", command, serverKey)
						fault = false
					}
				}
			}
		}
		if fault {
			// problem connecting to the distributed server...
			fmt.Printf("cluster: removing server '%s'\n", serverKey)
			// remove it...
			kvs.servers.Delete(serverKey)
		}
	}
}

func (kvs *KvServer) handleInternalChecking() {
	for {
		time.Sleep(5 * time.Second)
		_ = kvs.handleMessage(nil, &commandMessage{Command: "chk", Key: "", Value: ""})
	}
}

func (kvs *KvServer) Close() {

}

func (kvs *KvServer) WaitForShutdown() {
	<-kvs.shutdown
}

func (kvs *KvServer) Shutdown() {
	kvs.shutdown <- 1
}

func (kvs *KvServer) handleTcpAcceptance(listener net.Listener) {
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("server-tcp: error accepting connection: %s\n", err.Error())
			continue
		}
		if connection == nil {
			fmt.Printf("server-tcp: connection was nil\n")
			continue
		}
		go kvs.handleTcpConnection(connection)
	}
}

func (kvs *KvServer) handleTcpConnection(connection io.ReadWriteCloser) {
	defer func() { _ = connection.Close() }()

	parser, _ := parsing.NewParser(kvs.grammar)

	buffer := make([]byte, KvServerReadBufferSize)
	for {
		count, err := connection.Read(buffer)
		if err != nil {
			return
		}
		if count == 0 {
			continue
		}
		cont, err := kvs.handleReceivedBytes(connection, parser, buffer[:count])
		if !cont || err != nil {
			return
		}
	}
}

func (kvs *KvServer) handleReceivedBytes(connection io.Writer, parser *parsing.Parser, values []byte) (carryOn bool, e error) {
	for _, value := range values {
		cont, err := kvs.handleReceivedByte(connection, parser, value)
		if !cont || err != nil {
			return false, err
		}
	}
	return true, nil
}

func (kvs *KvServer) handleReceivedByte(connection io.Writer, parser *parsing.Parser, value byte) (carryOn bool, e error) {
	found, err := parser.Process(string(value))
	if err != nil {
		_, err := writeErr(connection)
		if err != nil {
			return false, err
		}
	}
	if found {
		cmd, arg1, arg2, err := parser.GetMessage()
		if err != nil {
			panic("server-tcp: something really vile has happened")
		}
		if !kvs.handleMessage(connection, &commandMessage{Command: cmd, Key: arg1, Value: arg2}) {
			return false, nil
		}
	}
	return true, nil
}

func writeErr(connection io.Writer) (n int, err error) {
	return connection.Write([]byte("err"))
}
