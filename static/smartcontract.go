package static

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/topic"

	zmq "github.com/pebbe/zmq4"
)

type (
	SmartcontractKey string
	Smartcontract    struct {
		// Body abi.ABI
		NetworkId           string
		Address             string
		AbiHash             string
		Txid                string
		Deployer            string
		StartingBlockNumber int
		exists              bool
	}
)

func CreateSmartcontractKey(networkId string, address string) SmartcontractKey {
	key := networkId + "." + address
	return SmartcontractKey(key)
}

func (c *Smartcontract) Key() (string, string) {
	return c.NetworkId, c.Address
}

func (c *Smartcontract) KeyString() SmartcontractKey {
	key := c.NetworkId + "." + c.Address
	return SmartcontractKey(key)
}

func (c *Smartcontract) SetExists(exists bool) {
	c.exists = exists
}

func (k *SmartcontractKey) Decompose() (string, string) {
	str := strings.Split(string(*k), ".")
	return str[0], str[1]
}

func NewSmartcontract(body map[string]interface{}) *Smartcontract {
	smartcontract := Smartcontract{}
	smartcontract.exists = false
	if body["network_id"] != nil {
		smartcontract.NetworkId = body["network_id"].(string)
	}
	if body["address"] != nil {
		smartcontract.Address = body["address"].(string)
	}
	if body["abi_hash"] != nil {
		smartcontract.AbiHash = body["abi_hash"].(string)
	}
	if body["txid"] != nil {
		smartcontract.Txid = body["txid"].(string)
	}
	// optional parameter
	if body["deployer"] != nil {
		smartcontract.Deployer = body["deployer"].(string)
	} else {
		smartcontract.Deployer = ""
	}
	if body["starting_block_number"] != nil {
		smartcontract.StartingBlockNumber = int(body["starting_block_number"].(float64))
	}

	return &smartcontract
}

func (smartcontract *Smartcontract) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = smartcontract.NetworkId
	i["address"] = smartcontract.Address
	i["abi_hash"] = smartcontract.AbiHash
	i["txid"] = smartcontract.Txid
	i["starting_block_number"] = smartcontract.StartingBlockNumber
	i["deployer"] = smartcontract.Deployer

	return i
}

func (smartcontract *Smartcontract) ToString() string {
	s := smartcontract.ToJSON()
	byt, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(byt)
}

// Returns list of smartcontracts by topic filter in remote Static service
func FilterSmartcontracts(socket *zmq.Socket, tf *topic.TopicFilter) []*Smartcontract {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_filter",
		Param: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	fmt.Println("Sending message to STATIC server to get smartcontracts. The mesage sent to server")
	fmt.Println(request.ToString())
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for smartcontracts getting from static controller", err.Error())
		return []*Smartcontract{}
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return []*Smartcontract{}
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse smartcontracts reply", err.Error())
		return []*Smartcontract{}
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return []*Smartcontract{}
	}

	rawSmartcontracts := reply.Params["smartcontracts"].([]map[string]interface{})
	var smartcontracts []*Smartcontract = make([]*Smartcontract, len(rawSmartcontracts))
	for i, rawSmartcontract := range rawSmartcontracts {
		smartcontracts[i] = NewSmartcontract(rawSmartcontract)
	}

	return smartcontracts
}

func FilterSmartcontractKeys(socket *zmq.Socket, tf *topic.TopicFilter) []SmartcontractKey {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_key_filter",
		Param: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	fmt.Println("Sending message to STATIC server to get smartcontract keys. The mesage sent to server")
	fmt.Println(request.ToString())
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for smartcontract keys getting from static controller", err.Error())
		return []SmartcontractKey{}
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return []SmartcontractKey{}
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse reply", err.Error())
		return []SmartcontractKey{}
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return []SmartcontractKey{}
	}

	rawKeys := reply.Params["smartcontract_keys"].([]interface{})
	var keys []SmartcontractKey = make([]SmartcontractKey, len(rawKeys))
	for i, rawKey := range rawKeys {
		keys[i] = SmartcontractKey(rawKey.(string))
	}

	return keys
}

func GetRemoteSmartcontract(socket *zmq.Socket, networkId string, address string) (*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get",
		Param: map[string]interface{}{
			"network_id": networkId,
			"address":    address,
		},
	}
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for smartcontract getting from static controller", err.Error())
		return nil, fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return nil, fmt.Errorf("receiving: %w", err)
	}

	fmt.Println(r)
	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse smartcontract reply", err.Error())
		return nil, fmt.Errorf("spaghetti block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure", reply.Message)
		return nil, fmt.Errorf("spaghetti block reply status is not ok: %s", reply.Message)
	}

	returnedSmartcontract := reply.Params["smartcontract"].(map[string]interface{})
	return NewSmartcontract(returnedSmartcontract), nil
}
