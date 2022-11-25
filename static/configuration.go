package static

import (
	"encoding/json"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/topic"
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

func (c *Configuration) SetId(id uint) {
	c.exists = true
	c.id = id
}

func (c *Configuration) SetAddress(address string) {
	c.Address = address
}

func (c *Configuration) Exists() bool { return c.exists }

// Creates a new static.Configuration class based on the given data
func NewConfiguration(body map[string]interface{}) *Configuration {
	conf := Configuration{
		Organization: body["o"].(string),
		Project:      body["p"].(string),
		NetworkId:    body["n"].(string),
		Group:        body["g"].(string),
		Name:         body["s"].(string),
	}
	address := ""
	if body["address"] != nil {
		address = body["address"].(string)
	}

	id := uint(0)
	if body["id"] != nil {
		id = uint(body["id"].(float64))
	}
	conf.id = id
	conf.exists = true
	conf.Address = address

	return &conf
}

func (c *Configuration) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"s":       c.Name,
		"n":       c.NetworkId,
		"g":       c.Group,
		"o":       c.Organization,
		"p":       c.Project,
		"address": c.Address,
	}
}

func (c *Configuration) ToString() string {
	interfaces := c.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return ""
	}

	return string(byt)
}

// get configuration from SDS Static by the configuration topic
func RemoteConfiguration(socket *remote.Socket, t *topic.Topic) (*Configuration, *Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "configuration_get",
		Param:   t.ToJSON(),
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, nil, err
	}

	returnedConfig := params["configuration"].(map[string]interface{})
	returnedSmartcontract := params["smartcontract"].(map[string]interface{})
	return NewConfiguration(returnedConfig), NewSmartcontract(returnedSmartcontract), nil
}

func RemoteConfigurationRegister(socket *remote.Socket, conf *Configuration) error {
	// Send hello.
	request := message.Request{
		Command: "configuration_register",
		Param:   conf.ToJSON(),
	}

	_, err := socket.RequestRemoteService(&request)
	return err
}
