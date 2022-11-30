package remote

import (
	"errors"
	"time"

	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

func (socket *Socket) SetSubscribeFilter(topic string) error {
	socketType, err := socket.socket.GetType()
	if err != nil {
		return err
	}
	if socketType != zmq.SUB {
		return errors.New("the socket is not a Broadcast. Can not call subscribe")
	}

	return socket.socket.SetSubscribe(topic)
}

func (socket *Socket) Subscribe(channel chan message.Reply, timeOut time.Duration) {
	socketType, err := socket.socket.GetType()
	if err != nil {
		channel <- message.Fail(err.Error())
		return
	}
	if socketType != zmq.SUB {
		channel <- message.Fail("the socket is not a Broadcast. Can not call subscribe")
		return
	}

	fetched := false
	time.AfterFunc(timeOut, func() {
		if !fetched {
			channel <- message.Fail("timeout for a message.\nPlease make sure SDS Spaghetti is running!")
		}
	})

	for {
		msgRaw, err := socket.socket.RecvMessage(0)

		fetched = true

		if err != nil {
			channel <- message.Fail(err.Error())
			continue
		}
		msgStr := append(msgRaw...)

		if len(msgStr) == 0 {
			channel <- message.Fail("Reconnect to the socket, it returned invalid message")
			continue
		}
		broadcast, err := message.ParseBroadcast(msg_raw)
		if err != nil {
			channel <- message.Fail("Error when parsing message: " + err.Error())
			continue
		}

		channel <- broadcast.Reply()
	}
}
