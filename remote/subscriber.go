package remote

import (
	"errors"
	"time"

	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

// The Socket if its a Subscriber applies a filter to listen certain data from the Broadcaster.
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

// Subscribe to the SDS Broadcast.
// The function is intended to be called as a gouritine.
//
// When a new message arrives, the method will send to the channel.
func (socket *Socket) Subscribe(channel chan message.Reply, time_out time.Duration) {
	socketType, err := socket.socket.GetType()
	if err != nil {
		channel <- message.Fail("failed to check the socket type. the socket error: " + err.Error())
		return
	}
	if socketType != zmq.SUB {
		channel <- message.Fail("the socket is not a Broadcast. Can not call subscribe")
		return
	}

	timer := time.AfterFunc(time_out, func() {
		channel <- message.Fail("timeout")
	})

	for {
		msgRaw, err := socket.socket.RecvMessage(0)

		timer.Reset(time_out)

		if err != nil {
			channel <- message.Fail(err.Error())
			continue
		}

		broadcast, err := message.ParseBroadcast(msgRaw)
		if err != nil {
			channel <- message.Fail("Error when parsing message: " + err.Error())
			continue
		}

		channel <- broadcast.Reply()
	}
}
