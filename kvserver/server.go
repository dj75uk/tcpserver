package kvserver

import (
	"errors"
	"fmt"
	"kvsapp/kvstore"
	"kvsapp/parsing"
	"net"
)

const KvServerReadBufferSize int = 256

type KvServer struct {
	port                    int
	store                   kvstore.KvStore
	verbsAndParameterCounts map[string]uint16
}

func NewKvServer(port int, store *kvstore.KvStore) (*KvServer, error) {
	if store == nil {
		return nil, errors.New("parameter 'store' must not be nil")
	}
	verbs := map[string]uint16{
		"put": 2,
		"get": 1,
		"del": 1,
		"bye": 0,
	}
	return &KvServer{
		port:                    port,
		store:                   *store,
		verbsAndParameterCounts: verbs,
	}, nil
}

func (kvs *KvServer) Open() error {

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", kvs.port))
	if err != nil {
		return err
	}
	fmt.Printf("server: listening on port %d\n", kvs.port)
	go kvs.handleAcceptance(listener)
	return nil
}

func (kvs *KvServer) Close() {

}

func (kvs *KvServer) handleAcceptance(listener net.Listener) {

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("server: error accepting connection: %s\n", err.Error())
			continue
		}
		if connection == nil {
			fmt.Printf("server: connection was nil\n")
			continue
		}
		go kvs.handleConnection(connection)
	}
}

func (kvs *KvServer) handleConnection(connection net.Conn) {
	defer func() { _ = connection.Close() }()

	parser := parsing.NewParser2()

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

func (kvs *KvServer) handleReceivedBytes(connection net.Conn, parser *parsing.Parser2, values []byte) (carryOn bool, e error) {
	for _, value := range values {
		cont, err := kvs.handleReceivedByte(connection, parser, value)
		if !cont || err != nil {
			return false, err
		}
	}
	return true, nil
}

func (kvs *KvServer) handleReceivedByte(connection net.Conn, parser *parsing.Parser2, value byte) (carryOn bool, e error) {
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
			panic("something really vile has happened")
		}
		if !kvs.handleMessage(connection, &parsing.Msg{Command: cmd, Key: arg1, Value: arg2}) {
			return false, nil
		}
	}
	return true, nil
}

func writeErr(connection net.Conn) (n int, err error) {
	return connection.Write([]byte("err"))
}

func (kvs *KvServer) handleMessage(connection net.Conn, message *parsing.Msg) (carryOn bool) {
	if message == nil {
		return false
	}

	fmt.Printf("processing...\n")
	fmt.Printf("- msg.Command: %s\n", message.Command)
	fmt.Printf("- msg.Key:     %s\n", message.Key)
	fmt.Printf("- msg.Value:   %s\n", message.Value)

	responseToWrite := "err"

	switch message.Command {
	case "bye":
		return false
	case "get":
		if result, err := kvs.store.Get(message.Key); err != nil {
			responseToWrite = "nil"
		} else {
			if bytesToWrite, err := parsing.NewParser().CreateData("val", result, ""); err == nil {
				responseToWrite = string(bytesToWrite)
			}
		}
	case "put":
		if _, err := kvs.store.Upsert(message.Key, message.Value); err == nil {
			responseToWrite = "ack"
		}
	case "del":
		if _, err := kvs.store.Delete(message.Key); err == nil {
			responseToWrite = "ack"
		}
	default:
		fmt.Printf("server: sending err (msg.Command == %s)\n", message.Command)
	}
	if len(responseToWrite) > 0 {
		_, err := connection.Write([]byte(responseToWrite))
		if err != nil {
			return false
		}
	}
	return true
}
