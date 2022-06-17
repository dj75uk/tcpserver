package kvserver

import (
	"fmt"
	"io"
	"kvsapp/parsing"
	"strconv"
)

func getHandlers() map[string]func(kvs *KvServer, key string, value string) string {
	return map[string]func(kvs *KvServer, key string, value string) string{
		"nop": handleNop,
		"chk": handleChk,
		"put": handlePut,
		"get": handleGet,
		"del": handleDel,
		"hed": handleHed,
		"spt": handleSpt,
		"sdl": handleSdl,
		"hst": handleHst,
	}
}

func (kvs *KvServer) handleMessage(connection io.Writer, message *commandMessage) (carryOn bool) {
	if message == nil {
		return false
	}

	fmt.Printf("server: handling '%s' command\n", message.Command)

	responseToWrite := "err"

	switch message.Command {

	case "bye":
		return false

	case "die":
		_ = handleDie(kvs, message.Key, message.Value)
		return false

	default:
		if handler, exists := kvs.handlers[message.Command]; exists {
			responseToWrite = handler(kvs, message.Key, message.Value)
		} else {
			fmt.Println("server: unknown command")
		}
	}
	if connection != nil && len(responseToWrite) > 0 {
		if responseToWrite == "err" {
			fmt.Printf("server: returning 'err' from: {%s}{%s}{%s}\n", message.Command, message.Key, message.Value)
		}
		_, err := connection.Write([]byte(responseToWrite))
		if err != nil {
			return false
		}
	}
	return true
}

func handleNop(kvs *KvServer, key string, value string) string {
	return "ack"
}

func handleDie(kvs *KvServer, key string, value string) string {
	kvs.Shutdown()
	return "ack"
}

func handleChk(kvs *KvServer, key string, value string) string {
	kvs.sendToAllOthers("nop", "", "")
	return "ack"
}

func handleHst(kvs *KvServer, key string, value string) string {
	// udp broadcast message...
	kvs.servers.Upsert(key, value)
	for _, serverKey := range kvs.servers.ListKeys() {
		fmt.Printf("cluster: server '%s' is currently known\n", serverKey)
	}
	return ""
}

func handlePut(kvs *KvServer, key string, value string) string {
	if _, err := kvs.store.Upsert(key, value); err == nil {
		kvs.sendToAllOthers("spt", key, value)
		return "ack"
	}
	return "err"
}

func handleSpt(kvs *KvServer, key string, value string) string {
	kvs.store.Upsert(key, value)
	return "ack"
}

func handleSdl(kvs *KvServer, key string, value string) string {
	kvs.store.Delete(key)
	return "ack"
}

func handleDel(kvs *KvServer, key string, value string) string {
	if _, err := kvs.store.Delete(key); err == nil {
		kvs.sendToAllOthers("sdl", key, "")
		return "ack"
	}
	return "err"
}

func handleGet(kvs *KvServer, key string, value string) string {
	if result, err := kvs.store.Get(key); err != nil {
		return "nil"
	} else {
		if bytesToWrite, err := parsing.CreateData("val", result, ""); err == nil {
			return string(bytesToWrite)
		}
	}
	return "err"
}

func handleHed(kvs *KvServer, key string, value string) string {
	if result, err := kvs.store.Get(key); err != nil {
		return "nil"
	} else {
		desiredLength, err := strconv.Atoi(value)
		if err == nil && desiredLength >= 0 {
			if desiredLength > 0 {
				result = result[0:desiredLength]
			}
			if bytesToWrite, err := parsing.CreateData("val", result, ""); err == nil {
				return string(bytesToWrite)
			}
		}
	}
	return "err"
}
