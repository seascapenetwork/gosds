package message

import (
	"encoding/json"
	"fmt"
	"strings"
)

// The PUB broadcasters broadcasting messages.
type Broadcast struct {
	Topic string
	reply Reply
}

// Convert to format understood by the protocol
func (b *Broadcast) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"topic": b.Topic,
		"reply": b.reply.ToJSON(),
	}
}

// broadcast as a string
func (b *Broadcast) ToString() string {
	return string(b.ToBytes())
}

// broadcast as a sequence of bytes
func (reply *Broadcast) ToBytes() []byte {
	interfaces := reply.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return []byte{}
	}

	return byt
}

// create a new broadcast
func NewBroadcast(topic string, reply Reply) Broadcast {
	return Broadcast{
		Topic: topic,
		reply: reply,
	}
}

// broadcast's actual data for the subscriber
func (b *Broadcast) Reply() Reply {
	return b.reply
}

// Is OK
func (r *Broadcast) IsOK() bool { return r.reply.IsOK() }

// parse the zeromq messages into a broadcast
func ParseBroadcast(msgs []string) (Broadcast, error) {
	msg := ""
	for _, v := range msgs {
		msg += v
	}
	i := strings.Index(msg, "{")
	topic := msg[:i]
	broadcastRaw := msg[i:]

	var dat map[string]interface{}

	if err := json.Unmarshal([]byte(broadcastRaw), &dat); err != nil {
		return Broadcast{}, err
	}

	if dat["reply"] == nil {
		return Broadcast{}, fmt.Errorf("no 'reply' parameter")
	}

	reply, replyErr := ParseJsonReply(dat["reply"].(map[string]interface{}))
	if replyErr != nil {
		return Broadcast{}, replyErr
	}

	return Broadcast{Topic: topic, reply: reply}, nil
}
