package message

import (
	"encoding/json"
)

// The SDS Service will get a request
type Request struct {
	Command string
	Param   map[string]interface{}
}

// Convert to JSON
func (reply *Request) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"command": reply.Command,
		"params":  reply.Param,
	}
}

// Convert request to the sequence of bytes
func (reply *Request) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}
	}

	return byt
}

// Convert request to the string
func (reply *Request) ToString() string {
	return string(reply.ToBytes())
}

// Messages from zmq concatenated
func ToString(msgs []string) string {
	msg := ""
	for _, v := range msgs {
		msg += v
	}
	return msg
}

// Parse the messages from zeromq into the Request
func ParseRequest(msgs []string) (Request, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		return Request{}, err
	}

	request := Request{
		Command: dat["command"].(string),
		Param:   dat["params"].(map[string]interface{}),
	}

	return request, nil
}
