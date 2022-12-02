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
	if body["pre_deploy_block_number"] != nil {
		smartcontract.PreDeployBlockNumber = int(body["pre_deploy_block_number"].(float64))
	}
	if body["pre_deploy_block_timestamp"] != nil {
		smartcontract.PreDeployBlockTimestamp = int(body["pre_deploy_block_timestamp"].(float64))
	}

	return &smartcontract
}

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
