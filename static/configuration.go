package static

import (
	"encoding/json"

	"github.com/blocklords/gosds/topic"

	zmq "github.com/pebbe/zmq4"
)

type Configuration struct {
	Organization string
	Project      string
	NetworkId    string
	Group        string
	Name         string
	Address      string
	id           uint
	exists       bool
}

func (c *Configuration) Exists() bool { return c.exists }

// Load Configuration from sds-static controller
func RemoteConfigByTopic(sock *zmq.Socket, t *topic.Topic) Configuration {
	return Configuration{}
}

func (c *Configuration) GetSmartcontract(networkId string, group string, name string) *Smartcontract {
	_, ok := c.config[networkId]
	if !ok {
		return nil
	}

	networkConfig := c.config[networkId]
	return networkConfig.GetSmartcontract(group, name)
}

func (c *Configuration) ExistSmartcontract(networkId string, group string, name string) bool {
	networkConfig, ok := c.config[networkId]
	if !ok {
		return false
	}

	groupConfig, ok2 := networkConfig.config[group]
	if !ok2 {
		return false
	}

	return groupConfig.SmartcontractInGroup(name)
}

func (c *Configuration) ToJSON() map[string]interface{} {
	obj := make(map[string]interface{}, len(c.config))

	for k, v := range c.config {
		obj[k] = v.ToJSON()
	}

	return obj
}

func (c *Configuration) ToString() string {
	interfaces := c.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return ""
	}

	return string(byt)
}
