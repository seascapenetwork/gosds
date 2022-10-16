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
		name                string
		address             string
		abiHash             string
		txid                string
		deployer            string
		startingBlockNumber int
		exists              bool
	}
)

func CreateSmartcontractKey(networkId string, address string) SmartcontractKey {
	key := networkId + "." + address
	return SmartcontractKey(key)
}

func NewSmartcontract(body map[string]interface{}) *Smartcontract {
	smartcontract := Smartcontract{}
	smartcontract.exists = false
	if body["name"] != nil {
		smartcontract.name = body["name"].(string)
	}
	if body["address"] != nil {
		smartcontract.address = body["address"].(string)
	}
	if body["abi_hash"] != nil {
		smartcontract.abiHash = body["abi_hash"].(string)
	}
	if body["txid"] != nil {
		smartcontract.txid = body["txid"].(string)
	}
	// optional parameter
	if body["deployer"] != nil {
		smartcontract.deployer = body["deployer"].(string)
	} else {
		smartcontract.deployer = ""
	}
	if body["starting_block_number"] != nil {
		smartcontract.startingBlockNumber = int(body["starting_block_number"].(float64))
	}

	return &smartcontract
}

func (smartcontract *Smartcontract) Address() string {
	return smartcontract.address
}

func (smartcontract *Smartcontract) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["name"] = smartcontract.name
	i["address"] = smartcontract.address
	i["abi_hash"] = smartcontract.abiHash
	i["txid"] = smartcontract.txid
	i["starting_block_number"] = smartcontract.startingBlockNumber
	i["deployer"] = smartcontract.deployer

	return i
}

func (smartcontract *Smartcontract) Name() string { return smartcontract.name }

func (smartcontract *Smartcontract) ToString() string {
	s := smartcontract.ToJSON()
	byt, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(byt)
}

// Returns list of smartcontracts by topic filter.
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
