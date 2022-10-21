package subscriber

import (
	"time"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/topic"

	"github.com/blocklords/gosds/sdk/remote"

	zmq "github.com/pebbe/zmq4"
)

type Subscriber struct {
	Host    string // SDS Gateway host
	Sub     string // SDS Publisher host
	Address string // Account address granted for reading
}

func NewSubscriber(host string, sub string, address string) *Subscriber {
	return &Subscriber{Host: host, Sub: sub, Address: address}
}

func (s *Subscriber) subscribe(t *topic.TopicFilter) {
	// preparing the subscriber so that we catch the first message if it was send
	// by publisher.
	time.Sleep(time.Millisecond * time.Duration(100))

	sender, senderErr := remote.NewPair(s.Address, true)
	if senderErr != nil {
		panic(senderErr)
	}

	msg := message.Request{
		Command: "subscribe",
		Param: map[string]interface{}{
			"topic_filter": t.ToJSON(),
			"subscriber":   s.Address,
		},
	}

	subscribed := remote.ReqReply(s.Host, msg)
	if !subscribed.IsOK() {
		sender.Send(subscribed.ToString(), 0)
		return
	}

	go s.heartbeat(sender)
}

func (s *Subscriber) heartbeat(timeout *zmq.Socket) {
	for {
		heartbeatReply := Heartbeat(s)
		if !heartbeatReply.IsOK() {
			timeout.Send(heartbeatReply.ToString(), 0)
			break
		}

		time.Sleep(time.Second * time.Duration(2))
	}
}

func (s *Subscriber) loop(sub *zmq.Socket, read *zmq.Socket, ch chan message.Broadcast) {
	poller := zmq.NewPoller()
	poller.Add(sub, zmq.POLLIN)
	poller.Add(read, zmq.POLLIN)

LOOP:
	for {
		sockets, _ := poller.Poll(-1)
		for _, socket := range sockets {
			switch sock := socket.Socket; sock {
			case sub:
				msg_raw, _ := sock.RecvMessage(0)
				b, err := message.ParseBroadcast(msg_raw)
				if err != nil {
					ch <- message.NewBroadcast(s.Address, message.Reply{Status: "fail", Message: "Error when parsing message " + err.Error()})
					break LOOP
				}

				//  Send results to sink
				ch <- b

				if !b.IsOK() {
					break LOOP //  Exit loop
				}
			case read:
				ch <- message.NewBroadcast(s.Address, message.Reply{Status: "fail", Message: "Stop the loop"})
				//  Any controller command acts as 'KILL'
				break LOOP //  Exit loop
			}
		}
	}
}

func (s *Subscriber) Listen(t *topic.TopicFilter) (message.Reply, chan message.Broadcast) {
	go s.subscribe(t)

	// Run the listener
	sub, err := remote.NewSub(s.Sub, s.Address)
	if err != nil {
		return message.Fail("Failed to establish a connection with SDS Publisher: " + err.Error()), nil
	}

	// Run heartbeat, subscription status reader
	read, readErr := remote.NewPair(s.Address, false)
	if readErr != nil {
		return message.Fail("Internal SDK Error: " + readErr.Error()), nil
	}

	// now create a heartbeat timer
	ch := make(chan message.Broadcast)

	go s.loop(sub, read, ch)

	return message.Reply{Status: "OK"}, ch
}
