/*Categorized log containing log name and output parameters*/
package categorizer

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/spaghetti"
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
		LogIndex:  uint(blob["log_index"].(float64)),
		Address:   blob["address"].(string),
		Log:       blob["log"].(string),
		Output:    blob["output"].(map[string]interface{}),
	}
}

func RemoteLogs(socket *remote.Socket, keys []string) ([]*Log, error) {
	request := message.Request{
		Command: "log_get_all",
		Param: map[string]interface{}{
			"keys": keys,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	logRaws := params["logs"].([]interface{})
	logs := make([]*Log, len(logRaws))
	for i, raw := range logRaws {
		logs[i] = ParseLog(raw.(map[string]interface{}))
	}

	return logs, nil
}

// parse the raw event data from spaghetti using SDS Log
// parsing events using JSON abi is harder in golang, therefore we use javascript
// implementation called SDS Log.
func RemoteLogParse(socket *remote.Socket, networkId string, address string, data string, topics []string) (string, interface{}, error) {
	request := message.Request{
		Command: "parse",
		Param: map[string]interface{}{
			"network_id": networkId,
			"address":    address,
			"data":       data,
			"topics":     topics,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return "", nil, err
	}

	return params["name"].(string), params["args"], nil
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
