package subscriber

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/sdk/remote"
)

func Heartbeat(s *Subscriber) message.Reply {

	msg := message.Request{
		Command: "heartbeat",
		Param: map[string]interface{}{
			"subscriber": s.Address,
		},
	}

	return remote.ReqReply(s.Host, msg)
}
