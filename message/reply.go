package message

import (
	"encoding/json"
	"fmt"
)

// SDS Service returns the reply.
type Reply struct {
	Status  string
	Message string
	Params  map[string]interface{}
}

func Fail(err string) Reply {
	return Reply{Status: "fail", Message: err}
}

// Is SDS Service returned a successful reply
func (r *Reply) IsOK() bool { return r.Status == "OK" }

// Convert to JSON
func (reply *Reply) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"status":  reply.Status,
		"message": reply.Message,
		"params":  reply.Params,
	}
}

// Convert the reply to the string format
func (reply *Reply) ToString() string {
	return string(reply.ToBytes())
}

// Reply as a sequence of bytes
func (reply *Reply) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}
	}

	return byt
}

// Zeromq received raw strings converted to the reply.
func ParseReply(msgs []string) (Reply, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}
	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(msg), &dat); err != nil {
		return Reply{}, err
	}

	return ParseJsonReply(dat)
}

// Reply object from a json object.
func ParseJsonReply(dat map[string]interface{}) (Reply, error) {
	if dat["status"] == nil {
		return Reply{}, fmt.Errorf("no 'status' parameter")
	}

	replyMessage := ""
	if dat["message"] != nil {
		replyMessage = dat["message"].(string)
	}

	var params map[string]interface{}
	if dat["params"] != nil {
		params = dat["params"].(map[string]interface{})
	}

	reply := Reply{
		Status:  dat["status"].(string),
		Params:  params,
		Message: replyMessage,
	}

	return reply, nil
}
