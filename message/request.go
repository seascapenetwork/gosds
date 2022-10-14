package message

import (
	"encoding/json"
)

type Request struct {
	Command string
	Param   map[string]interface{}
}

func (reply *Request) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"command": reply.Command,
		"params":  reply.Param,
	}
}

func (reply *Request) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}
	}

	return byt
}

func (reply *Request) ToString() string {
	return string(reply.ToBytes())
}

func ToString(msgs []string) string {
	msg := ""
	for _, v := range msgs {
		msg += v
	}
	return msg
}

func ParseRequest(msgs []string) (Request, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		return Request{}, err
	}

	return Request{Command: dat["command"].(string), Param: dat["params"].(map[string]interface{})}, nil
}
