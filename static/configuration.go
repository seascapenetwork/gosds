package static

import (
	"encoding/json"
	"fmt"

	"github.com/blocklords/gosds/message"
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
func RemoteConfigByTopic(socket *zmq.Socket, t *topic.Topic) (*Configuration, error) {
	// Send hello.
	request := message.Request{
		Command: "configuration_get",
		Param: map[string]interface{}{
			"organization": t.Organization,
			"project":      t.Project,
			"group":        t.Group,
			"network_id":   t.NetworkId,
			"name":         t.Name,
		},
	}
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to get config from SDS-Static", err.Error())
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
		fmt.Println("Failed to parse abi reply", err.Error())
		return nil, fmt.Errorf("spaghetti block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure", reply.Messge)
		return nil, fmt.Errorf("spaghetti block reply status is not ok: %s", reply.Message)
	}

	returnedSmartcontract := reply.Params["configuration"].(map[string]interface{})
	return NewConfiguration(returnedSmartcontract), nil
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
