package message

import (
	"encoding/json"
	"strings"
)

type Broadcast struct {
	Topic string
	reply Reply
}

/* Convert to format understood by the protocol */
func (b *Broadcast) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"topic":   b.Topic,
		"message": b.reply.ToJSON(),
	}
}

func (b *Broadcast) ToString() string {
	return string(b.ToBytes())
}

func (reply *Broadcast) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}
	}

	return byt
}

func NewBroadcast(topic string, reply Reply) Broadcast {
	return Broadcast{
		Topic: topic,
		reply: reply,
	}
}

func (b *Broadcast) Reply() Reply {
	return b.reply
}

func (r *Broadcast) IsOK() bool { return r.reply.Status == "ok" || r.reply.Status == "OK" }

func ParseBroadcast(msgs []string) (Broadcast, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}
	i := strings.Index(msg, "{")
	topic := msg[:i]
	replyRaw := msg[i:]

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(replyRaw), &dat); err != nil {
		return Broadcast{}, err
	}

	replyMessage := ""
	if dat["message"] != nil && len(dat["message"].(string)) > 0 {
		replyMessage = dat["message"].(string)
	}

	reply := Reply{
		Status:  dat["status"].(string),
		Message: replyMessage,
		Params:  dat["params"].(map[string]interface{}),
	}

	return Broadcast{Topic: topic, reply: reply}, nil
}
