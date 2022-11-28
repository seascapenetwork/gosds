// The sdk/subscriber package is used to register for the smartcontracts
package subscriber

import (
	"fmt"
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
	timer   *time.Timer
}

func NewSubscriber(host string, sub string, address string) *Subscriber {
	return &Subscriber{Host: host, Sub: sub, Address: address}
}

func (s *Subscriber) subscribe(t *topic.TopicFilter, ch chan message.Reply) {
	// preparing the subscriber so that we catch the first message if it was send
	// by publisher.
	time.Sleep(time.Millisecond * time.Duration(100))

	msg := message.Request{
		Command: "subscribe",
		Param: map[string]interface{}{
			"topic_filter": t.ToJSON(),
			"subscriber":   s.Address,
		},
	}

	subscribed := remote.ReqReply(s.Host, msg)
	if !subscribed.IsOK() {
		ch <- subscribed
		return
	}

	s.timer = time.AfterFunc(time.Second*time.Duration(10), func() {
		ch <- message.Reply{Status: "fail", Message: "Server is not responding"}
	})

	go s.heartbeat(ch)
}

func (s *Subscriber) heartbeat(ch chan message.Reply) {
	for {
		s.timer.Reset(time.Second * time.Duration(10))

		heartbeatReply := Heartbeat(s)
		if !heartbeatReply.IsOK() {
			ch <- heartbeatReply
			break
		}

		time.Sleep(time.Second * time.Duration(2))
	}
}

// func (s *Subscriber) loop(sub *zmq.Socket, read *zmq.Socket, ch chan message.Broadcast) {
func (s *Subscriber) loop(sub *zmq.Socket, ch chan message.Broadcast) {
	//  Process messages from both sockets
	//  We prioritize traffic from the task ventilator

	for {
		msg_raw, err := sub.RecvMessage(0)
		if err != nil {
			fmt.Println("error in sub receive")
			fmt.Println(err)
		}
		if len(msg_raw) == 0 {
			break
		}

		b, err := message.ParseBroadcast(msg_raw)
		if err != nil {
			ch <- message.NewBroadcast(s.Address, message.Reply{Status: "fail", Message: "Error when parsing message " + err.Error()})
			break
		}

		//  Send results to sink
		ch <- b

		if !b.IsOK() {
			break //  Exit loop
		}
	}
}

func (s *Subscriber) Listen(t *topic.TopicFilter) (message.Reply, chan message.Broadcast, chan message.Reply) {
	hb := make(chan message.Reply)

	go s.subscribe(t, hb)

	// Run the listener
	sub, err := remote.NewSub(s.Sub, s.Address)
	if err != nil {
		return message.Fail("Failed to establish a connection with SDS Publisher: " + err.Error()), nil, nil
	}

	// now create a heartbeat timer
	ch := make(chan message.Broadcast)

	go s.loop(sub, ch)

	return message.Reply{Status: "OK", Message: "Successfully created a listener"}, ch, hb
}
