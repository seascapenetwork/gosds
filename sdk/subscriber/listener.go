// Listener is the one that listens the publishing
package subscriber

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/sdk/remote"
)

func Listen(s *Subscriber, m chan message.Reply) message.Reply {
	msg := message.Request{
		Command: "heartbeat",
		Param: map[string]interface{}{
			"subscriber": s.Address,
		},
	}

	return remote.ReqReply(s.Host, msg)
}
