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

func writeErr(connection net.Conn) (n int, err error) {
	return connection.Write([]byte("err"))
}

func (kvs *KvServer) handleConnection(connection net.Conn) {

	fmt.Printf("server: handling connection from %v\n", connection.RemoteAddr())

	defer func() { _ = connection.Close() }()
	parser := parsing.NewParser()
	accumulator := make([]byte, 0, KvServerReadBufferSize*8)
	for {
		buffer := make([]byte, KvServerReadBufferSize)
		count, err := connection.Read(buffer)
		if err != nil {
			return
		}
		if count == 0 {
			continue
		}

		fmt.Printf("received...\n")
		fmt.Printf("- count:  %d\n", count)
		fmt.Printf("- err:    %v\n", err)
		fmt.Printf("- buffer: [%s]\n", string(buffer))
		fmt.Printf("\n")

		fmt.Printf("appending accumulator...\n")
		fmt.Printf("- previous: %s\n", string(accumulator))
		fmt.Printf("- count:    %d\n", count)
		accumulator = append(accumulator, buffer[0:count]...)
		fmt.Printf("- current:  %s\n", string(accumulator))
		fmt.Printf("\n")

		if len(accumulator) < 3 {
			continue
		}

		fmt.Printf("checking command...\n")
		cmd := string(accumulator[0:3])
		expectedParameters, exists := kvs.verbsAndParameterCounts[cmd]
		fmt.Printf("- cmd:                %s\n", cmd)
		fmt.Printf("- expectedParameters: %d\n", expectedParameters)
		fmt.Printf("- exists:             %v\n", exists)
		fmt.Printf("\n")

		if !exists {
			fmt.Printf("server: cmd: %s unknown!\n", cmd)
			fmt.Printf("decapitating accumulator...\n")
			accumulator = accumulator[3:] // get rid of the crap
			_, err := writeErr(connection)
			if err != nil {
				return
			}
			continue
		}

		fmt.Printf("parsing...\n")
		fmt.Printf("- accumulator:        %s\n", string(accumulator[0:30]))
		fmt.Printf("- len(accumulator):   %d\n", len(accumulator))
		fmt.Printf("- expectedParameters: %d\n", expectedParameters)
		message, chop, err := parser.Parse(accumulator, expectedParameters)
		fmt.Printf("- msg                 %v\n", message)
		fmt.Printf("- chop:               %d\n", chop)
		fmt.Printf("- err:                %v\n", err)
		fmt.Printf("\n")

		if chop > 0 {
			// remove the processed bytes from the accumulator...
			fmt.Printf("decapitating accumulator...\n")
			fmt.Printf("- previous: %s\n", string(accumulator))
			fmt.Printf("- chop:     %d\n", chop)
			accumulator = accumulator[chop:]
			fmt.Printf("- current:  %s\n", string(accumulator))
			fmt.Printf("\n")
		}
		if err != nil {
			// send "err" to client
			fmt.Printf("server: sending err (parser err == %s)\n", err.Error())
			_, err := writeErr(connection)
			if err != nil {
				return
			}
			continue
		}
		if message != nil {
			if !kvs.handleMessage(connection, message) {
				return
			}
		}
	}
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
