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
	// network id + "." + address
	SmartcontractKey string
	Smartcontract    struct {
		// Body abi.ABI
		NetworkId           string
		Address             string
		AbiHash             string
		Txid                string
		Deployer            string
		StartingBlockNumber int
		StartingTimestamp   int
		exists              bool
	}
	// smartcontract key => topic string
	FilteredSmartcontractKeys map[SmartcontractKey]string
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
	if body["starting_timestamp"] != nil {
		smartcontract.StartingTimestamp = int(body["starting_timestamp"].(float64))
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
	i["starting_timestamp"] = smartcontract.StartingTimestamp
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
// also the topic path of the smartcontract
func FilterSmartcontracts(socket *zmq.Socket, tf *topic.TopicFilter) ([]*Smartcontract, []string) {
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
		return []*Smartcontract{}, nil
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return []*Smartcontract{}, nil
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse smartcontracts reply", err.Error())
		return []*Smartcontract{}, nil
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return []*Smartcontract{}, nil
	}

	rawSmartcontracts := reply.Params["smartcontracts"].([]interface{})
	rawTopics := reply.Params["smartcontracts"].([]interface{})
	var smartcontracts []*Smartcontract = make([]*Smartcontract, len(rawSmartcontracts))
	var topicStrings []string = make([]string, len(rawSmartcontracts))
	for i, rawSmartcontract := range rawSmartcontracts {
		smartcontracts[i] = NewSmartcontract(rawSmartcontract.(map[string]interface{}))
		topicStrings[i] = rawTopics[i].(string)
	}

	return smartcontracts, topicStrings
}

func FilterSmartcontractKeys(socket *zmq.Socket, tf *topic.TopicFilter) FilteredSmartcontractKeys {
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
		return nil
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return nil
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse reply", err.Error())
		return nil
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return nil
	}

	rawKeys := reply.Params["smartcontract_keys"].(map[string]interface{})
	var keys FilteredSmartcontractKeys = make(FilteredSmartcontractKeys, len(rawKeys))
	for key, topicString := range rawKeys {
		keys[SmartcontractKey(key)] = topicString.(string)
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
