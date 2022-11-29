package message

import (
	"encoding/json"
	"errors"
)

// The SDS Service will get a request
type Request struct {
	Command string
	Param   map[string]interface{}
}

// Convert to JSON
func (request *Request) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"command": request.Command,
		"params":  request.Param,
	}
}

// Extracts the casted type
func (request *Request) ParameterUint64(name string) (uint64, error) {
	raw, exists := request.Param[name]
	if !exists {
		return 0, errors.New("missing '" + name + "' parameter in the Request")
	}
	value, ok := raw.(float64)
	if !ok {
		return 0, errors.New("expected number type for '" + name + "' parameter")
	}

	return uint64(value), nil
}

func (request *Request) ParameterString(name string) (string, error) {
	raw, exists := request.Param[name]
	if !exists {
		return "", errors.New("missing '" + name + "' parameter in the Request")
	}
	value, ok := raw.(string)
	if !ok {
		return "", errors.New("expected string type for '" + name + "' parameter")
	}

	return value, nil
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
