package static

import (
	"encoding/json"
	"errors"
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
		Parameters: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, nil, err
	}

	raw_smartcontracts, err := message.GetMapList(params, "smartcontracts")
	if err != nil {
		return nil, nil, err
	}
	topic_strings, err := message.GetStringList(params, "topics")
	if err != nil {
		return nil, nil, err
	}
	if len(raw_smartcontracts) != len(topic_strings) {
		return nil, nil, errors.New("the returned amount of topic strings mismatch with smartcontracts")
	}
	var smartcontracts []*Smartcontract = make([]*Smartcontract, len(raw_smartcontracts))
	for i, raw_smartcontract := range raw_smartcontracts {
		smartcontract, err := NewSmartcontract(raw_smartcontract)
		if err != nil {
			return nil, nil, err
		}
		smartcontracts[i] = smartcontract
	}

	return smartcontracts, topic_strings, nil
}

// returns list of smartcontract keys by topic filter
func RemoteSmartcontractKeys(socket *remote.Socket, tf *topic.TopicFilter) (FilteredSmartcontractKeys, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_key_filter",
		Parameters: map[string]interface{}{
			"topic_filter": tf.ToJSON(),
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	raw_keys, err := message.GetMap(params, "smartcontract_keys")
	if err != nil {
		return nil, err
	}
	var keys FilteredSmartcontractKeys = make(FilteredSmartcontractKeys, len(raw_keys))
	for key, raw_value := range raw_keys {
		topic_string, ok := raw_value.(string)
		if !ok {
			return nil, errors.New("one of the topic strings is not in the string format")
		}
		keys[SmartcontractKey(key)] = topic_string
	}

	return keys, nil
}

// returns smartcontract by smartcontract key from SDS Static
func RemoteSmartcontract(socket *remote.Socket, network_id string, address string) (*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get",
		Parameters: map[string]interface{}{
			"network_id": network_id,
			"address":    address,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	raw_smartcontract, err := message.GetMap(params, "smartcontract")
	if err != nil {
		return nil, err
	}
	return NewSmartcontract(raw_smartcontract)
}

func RemoteSmartcontractRegister(socket *remote.Socket, s *Smartcontract) (string, error) {
	// Send hello.
	request := message.Request{
		Command:    "smartcontract_register",
		Parameters: s.ToJSON(),
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return "", err
	}

	return message.GetString(params, "address")
}
