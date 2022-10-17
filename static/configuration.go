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
		Organization: body["organization"].(string),
		Project:      body["project"].(string),
		NetworkId:    body["network_id"].(string),
		Group:        body["group"].(string),
		Name:         body["name"].(string),
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

// Load Configuration from sds-static controller
func RemoteConfigByTopic(sock *zmq.Socket, t *topic.Topic) Configuration {
	return Configuration{}
}

func (c *Configuration) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"name":         c.Name,
		"network_id":   c.NetworkId,
		"group":        c.Group,
		"organization": c.Organization,
		"project":      c.Project,
		"address":      c.Address,
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
