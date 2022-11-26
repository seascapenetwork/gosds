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
	Nonce   uint
	Address string
}

func Fail(err string, address string, nonce uint) Reply {
	return Reply{Status: "fail", Message: err, Address: address, Nonce: nonce}
}

// Is SDS Service returned a successful reply
func (r *Reply) IsOK() bool { return r.Status == "OK" }

// Convert to JSON
func (reply *Reply) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"status":  reply.Status,
		"message": reply.Message,
		"params":  reply.Params,
		"address": reply.Address,
		"nonce":   reply.Nonce,
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

	address := dat["address"].(string)
	nonce := uint(dat["nonce"].(float64))

	var params map[string]interface{}
	if dat["params"] != nil {
		params = dat["params"].(map[string]interface{})
	}

	reply := Reply{
		Status:  dat["status"].(string),
		Params:  params,
		Message: replyMessage,
		Address: address,
		Nonce:   nonce,
	}

	return reply, nil
}
