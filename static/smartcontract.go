package static

import (
	"encoding/json"
	"fmt"

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

func (c *Smartcontract) SetExists(exists bool) {
	c.exists = exists
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
func FilterSmartcontracts(socket *zmq.Socket, tf *topic.TopicFilter) []Smartcontract {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_filter",
		Param: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	fmt.Println("Sending message to STATIC server to get abi. The mesage sent to server")
	fmt.Println(request.ToString())
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for abi getting from static controller", err.Error())
		return []Smartcontract{}
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller", err.Error())
		return []Smartcontract{}
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse abi reply", err.Error())
		return []Smartcontract{}
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure: ", reply.Message)
		return []Smartcontract{}
	}

	rawSmartcontracts := reply.Params["smartcontracts"].([]map[string]interface{})
	var smartcontracts []Smartcontract = make([]Smartcontract, len(rawSmartcontracts))
	for i, rawSmartcontract := range rawSmartcontracts {
		smartcontracts[i] = *NewSmartcontract(rawSmartcontract)
	}

	return smartcontracts
}
