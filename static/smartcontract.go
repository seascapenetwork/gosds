package static

import (
	"encoding/json"
	"strings"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/topic"
)

type (
	// network id + "." + address
	SmartcontractKey string
	Smartcontract    struct {
		// Body abi.ABI
		NetworkId               string
		Address                 string
		AbiHash                 string
		Txid                    string
		Deployer                string
		PreDeployBlockNumber    int
		PreDeployBlockTimestamp int
		exists                  bool
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

// The smartcontract parameters that composes the smartcontract key
// its the network id and the address
func (k *SmartcontractKey) Decompose() (string, string) {
	str := strings.Split(string(*k), ".")
	return str[0], str[1]
}

// Creates a new smartcontract from the JSON
func NewSmartcontract(parameters map[string]interface{}) (*Smartcontract, error) {
	network_id, err := message.GetString(parameters, "network_id")
	if err != nil {
		return nil, err
	}
	address, err := message.GetString(parameters, "address")
	if err != nil {
		return nil, err
	}
	abi_hash, err := message.GetString(parameters, "abi_hash")
	if err != nil {
		return nil, err
	}
	txid, err := message.GetString(parameters, "txid")
	if err != nil {
		return nil, err
	}
	// optional parameters
	deployer, err := message.GetString(parameters, "deployer")
	if err != nil {
		deployer = ""
	}
	pre_deploy_block_number, err := message.GetUint64(parameters, "pre_deploy_block_number")
	if err != nil {
		pre_deploy_block_number = 0
	}
	pre_deploy_block_timestamp, err := message.GetUint64(parameters, "pre_deploy_block_timestamp")
	if err != nil {
		return nil, err
	}

	smartcontract := Smartcontract{
		exists:                  false,
		NetworkId:               network_id,
		Address:                 address,
		AbiHash:                 abi_hash,
		Txid:                    txid,
		Deployer:                deployer,
		PreDeployBlockNumber:    int(pre_deploy_block_number),
		PreDeployBlockTimestamp: int(pre_deploy_block_timestamp),
	}
	return &smartcontract, nil
}

// JSON represantion of the static.Smartcontract
func (smartcontract *Smartcontract) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = smartcontract.NetworkId
	i["address"] = smartcontract.Address
	i["abi_hash"] = smartcontract.AbiHash
	i["txid"] = smartcontract.Txid
	i["pre_deploy_block_number"] = smartcontract.PreDeployBlockNumber
	i["pre_deploy_block_timestamp"] = smartcontract.PreDeployBlockTimestamp
	i["deployer"] = smartcontract.Deployer

	return i
}

// The JSON string represantion of the static.Smartcontract
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
func RemoteSmartcontracts(socket *remote.Socket, tf *topic.TopicFilter) ([]*Smartcontract, []string, error) {
	request := message.Request{
		Command: "smartcontract_filter",
		Param: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, nil, err
	}

	rawSmartcontracts := params["smartcontracts"].([]interface{})
	rawTopics := params["topics"].([]interface{})
	var smartcontracts []*Smartcontract = make([]*Smartcontract, len(rawSmartcontracts))
	var topicStrings []string = make([]string, len(rawSmartcontracts))
	for i, rawSmartcontract := range rawSmartcontracts {
		smartcontracts[i] = NewSmartcontract(rawSmartcontract.(map[string]interface{}))
		topicStrings[i] = rawTopics[i].(string)
	}

	return smartcontracts, topicStrings, nil
}

// returns list of smartcontract keys by topic filter
func RemoteSmartcontractKeys(socket *remote.Socket, tf *topic.TopicFilter) (FilteredSmartcontractKeys, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_key_filter",
		Param: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	rawKeys := params["smartcontract_keys"].(map[string]interface{})
	var keys FilteredSmartcontractKeys = make(FilteredSmartcontractKeys, len(rawKeys))
	for key, topicString := range rawKeys {
		keys[SmartcontractKey(key)] = topicString.(string)
	}

	return keys, nil
}

// returns smartcontract by smartcontract key from SDS Static
func RemoteSmartcontract(socket *remote.Socket, networkId string, address string) (*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get",
		Param: map[string]interface{}{
			"network_id": networkId,
			"address":    address,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	returnedSmartcontract := params["smartcontract"].(map[string]interface{})
	return NewSmartcontract(returnedSmartcontract), nil
}

func RemoteSmartcontractRegister(socket *remote.Socket, s *Smartcontract) (string, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_register",
		Param:   s.ToJSON(),
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return "", err
	}

	address := params["address"].(string)
	return address, nil
}
