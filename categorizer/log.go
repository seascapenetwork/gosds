/*Categorized log containing log name and output parameters*/
package categorizer

import (
	"fmt"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/spaghetti"

	zmq "github.com/pebbe/zmq4"
)

type Log struct {
	ID        uint64
	NetworkId string
	Txid      string
	LogIndex  uint
	Address   string
	Log       string
	Output    map[string]interface{}
}

func (b *Log) Key() string {
	return b.NetworkId + "." + b.Address
}

func (b *Log) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = b.NetworkId
	i["txid"] = b.Txid
	i["log_index"] = b.LogIndex
	i["address"] = b.Address
	i["log"] = b.Log
	i["output"] = b.Output
	return i
}

func ParseLog(blob map[string]interface{}) *Log {
	return &Log{
		NetworkId: blob["network_id"].(string),
		Txid:      blob["txid"].(string),
		LogIndex:  blob["log_index"].(uint),
		Address:   blob["address"].(string),
		Log:       blob["log"].(string),
		Output:    blob["output"].(map[string]interface{}),
	}
}

func RemoteLogs(socket *zmq.Socket, keys []string) ([]*Log, error) {
	request := message.Request{
		Command: "log_get_all",
		Param: map[string]interface{}{
			"keys": keys,
		},
	}
	fmt.Println("Sending message to SDS Log server to parse log. The mesage sent to server")
	fmt.Println(request.ToString())
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for smartcontracts getting from SDS Log", err.Error())
		return nil, err
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return nil, err
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse smartcontracts reply", err.Error())
		return nil, err
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return nil, err
	}

	logRaws := reply.Params["logs"].([]map[string]interface{})
	logs := make([]*Log, len(logRaws))
	for i, raw := range logRaws {
		logs[i] = ParseLog(raw)
	}

	return logs, nil
}

// parse the raw event data from spaghetti using SDS Log
// parsing events using JSON abi is harder in golang, therefore we use javascript
// implementation called SDS Log.
func RemoteParse(socket *zmq.Socket, networkId string, address string, data string, topics []string) (string, interface{}, error) {
	request := message.Request{
		Command: "parse",
		Param: map[string]interface{}{
			"network_id": networkId,
			"address":    address,
			"data":       data,
			"topics":     topics,
		},
	}
	fmt.Println("Sending message to SDS Log server to parse log. The mesage sent to server")
	fmt.Println(request.ToString())
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for smartcontracts getting from SDS Log", err.Error())
		return "", nil, err
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return "", nil, err
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse smartcontracts reply", err.Error())
		return "", nil, err
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return "", nil, err
	}

	return reply.Params["name"].(string), reply.Params["args"], nil
}

func NewLog(l spaghetti.Log, log string, output map[string]interface{}, c *Block) Log {
	return Log{
		NetworkId: c.NetworkID(),
		Address:   c.Address(),
		Txid:      l.TxId(),
		LogIndex:  l.LogIndex(),
		Log:       log,
		Output:    output,
	}
}
