package kvserver

import (
	"errors"
	"fmt"
	"io"
	"kvsapp/kvstore"
	"kvsapp/parsing"
	"net"
	"strconv"
)

const KvServerReadBufferSize int = 256

type KvServer struct {
	port     int
	store    kvstore.KvStore
	grammar  map[string]parsing.ParserGrammar
	shutdown chan int
}

type commandMessage struct {
	Command string
	Key     string
	Value   string
}

func NewKvServer(port int, store *kvstore.KvStore) (*KvServer, error) {
	if store == nil {
		return nil, errors.New("parameter 'store' must not be nil")
	}
	return &KvServer{
		port:  port,
		store: *store,
		grammar: map[string]parsing.ParserGrammar{
			"die": {ExpectedArguments: 0},
			"bye": {ExpectedArguments: 0},
			"get": {ExpectedArguments: 1},
			"del": {ExpectedArguments: 1},
			"put": {ExpectedArguments: 2},
			"hed": {ExpectedArguments: 2, Arg1LengthIsValue: false},
		},
		shutdown: make(chan int),
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

func (kvs *KvServer) WaitForShutdown() {
	<-kvs.shutdown
}

func (kvs *KvServer) Shutdown() {
	kvs.shutdown <- 1
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

func (kvs *KvServer) handleConnection(connection io.ReadWriteCloser) {
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
			panic("something really vile has happened")
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

func (kvs *KvServer) handleMessage(connection io.Writer, message *commandMessage) (carryOn bool) {
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
			if bytesToWrite, err := parsing.CreateData("val", result, ""); err == nil {
				responseToWrite = string(bytesToWrite)
			}
		}

	case "hed":
		if result, err := kvs.store.Get(message.Key); err != nil {
			responseToWrite = "nil"
		} else {
			desiredLength, err := strconv.Atoi(message.Value)
			if err == nil && desiredLength >= 0 {
				if bytesToWrite, err := parsing.CreateData("val", result, ""); err == nil {
					responseToWrite = string(bytesToWrite[0:desiredLength])
				}
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

	case "die":
		responseToWrite = "ack"
		kvs.Shutdown()
		return false

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
