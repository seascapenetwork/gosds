package writer

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/topic"

	"github.com/blocklords/gosds/sdk/remote"
)

type Writer struct {
	host    string // SDS Gateway host
	address string // Account address granted for reading
}

func NewWriter(host string, address string) *Writer {
	return &Writer{host: host, address: address}
}

func (r *Writer) Write(t topic.Topic, args map[string]interface{}) message.Reply {
	if t.Level() != topic.LEVEL_FULL {
		return message.Fail(`Topic should contain method name`)
	}

	msg := message.Request{
		Command: "smartcontract_write",
		Param: map[string]interface{}{
			"topic_string": t.ToString(topic.LEVEL_FULL),
			"arguments":    args,
			"address":      r.address,
		},
	}

	return remote.ReqReply(r.host, msg)
}
