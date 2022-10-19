package reader

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/topic"

	"github.com/blocklords/gosds/sdk/remote"
)

type Reader struct {
	host    string // SDS Gateway host
	address string // Account address granted for reading
}

func NewReader(host string, address string) *Reader {
	return &Reader{host: host, address: address}
}

func (r *Reader) Read(t *topic.Topic, args []interface{}) message.Reply {
	if t.Level() != topic.LEVEL_FULL {
		return message.Fail(`Topic should contain method name`)
	}

	msg := message.Request{
		Command: "smartcontract_read",
		Param: map[string]interface{}{
			"topic_string": t.ToString(topic.LEVEL_FULL),
			"arguments":    args,
			"address":      r.address,
		},
	}

	return remote.ReqReply(r.host, msg)
}
