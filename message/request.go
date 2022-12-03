package message

import (
	"encoding/json"
	"fmt"
)

// The SDS Service will accepts the Request message.
type Request struct {
	Command string
	Param   map[string]interface{}
}

// Convert Request to JSON
func (request *Request) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"command": request.Command,
		"params":  request.Param,
	}
}

// Request message as a  sequence of bytes
func (reply *Request) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		fmt.Println("error while converting json into bytes", err)
		return []byte{}
	}

	return byt
}

// Convert Request message to the string
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

	command, err := GetString(dat, "command")
	if err != nil {
		return Request{}, err
	}
	parameters, err := GetMap(dat, "params")
	if err != nil {
		return Request{}, err
	}

	request := Request{
		Command: command,
		Param:   parameters,
	}

	return request, nil
}
