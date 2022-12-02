package message

import (
	"encoding/json"
<<<<<<< HEAD
	"fmt"
	"time"
=======
	"errors"
>>>>>>> main
)

// SDS Service returns the reply. Anyone who sends a request to the SDS Service gets this message.
type Reply struct {
	Status  string
	Message string
	Params  map[string]interface{}
}

// Create a new Reply as a failure
// It accepts the error message that explains the reason of the failure.
func Fail(message string) Reply {
	return Reply{Status: "fail", Message: message, Params: map[string]interface{}{}}
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

// Zeromq received raw strings converted to the Reply message.
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

// Create 'Reply' message from a json.
func ParseJsonReply(dat map[string]interface{}) (Reply, error) {
	reply := Reply{}
	status, err := GetString(dat, "status")
	if err != nil {
		return reply, err
	}
	if status != "fail" && status != "OK" {
		return reply, errors.New("the 'status' of the reply can be either 'fail' or 'OK'.")
	} else {
		reply.Status = status
	}

<<<<<<< HEAD
	replyMessage := ""
	if dat["message"] != nil {
		fmt.Println("[", time.Now().Local(), "] got message: ", dat)
		jsonError, _ := json.Marshal(dat["message"])
		replyMessage = string(jsonError)
=======
	message, err := GetString(dat, "message")
	if err != nil {
		return reply, err
	} else {
		reply.Message = message
>>>>>>> main
	}

	parameters, err := GetMap(dat, "params")
	if err != nil {
		return reply, err
	} else {
		reply.Params = parameters
	}

	return reply, nil
}
