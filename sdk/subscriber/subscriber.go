// The sdk/subscriber package is used to register for the smartcontracts
package subscriber

import (
	"fmt"
	"time"

	"github.com/blocklords/gosds/message"
	sds_remote "github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/static"
	"github.com/blocklords/gosds/topic"

	"github.com/blocklords/gosds/sdk/db"
	"github.com/blocklords/gosds/sdk/remote"

	zmq "github.com/pebbe/zmq4"
)

type Subscriber struct {
	Address           string // Account address granted for reading
	timer             *time.Timer
	socket            *sds_remote.Socket
	db                *db.KVM                    // it also keeps the topic filter
	smartcontractKeys []*static.SmartcontractKey // list of smartcontract keys
}

func NewSubscriber(gatewaySocket *sds_remote.Socket, db *db.KVM, address string) *Subscriber {
	smartcontractKeys := make([]*static.SmartcontractKey, 0)

	return &Subscriber{
		Address:           address,
		socket:            gatewaySocket,
		db:                db,
		smartcontractKeys: smartcontractKeys,
	}
}

// The algorithm
// List of the smartcontracts by smartcontract filter
func (s *Subscriber) subscribe(ch chan message.Reply) error {
	// preparing the subscriber so that we catch the first message if it was send
	// by publisher.
	time.Sleep(time.Millisecond * time.Duration(100))

	request := message.Request{
		Command: "subscribe",
		Param: map[string]interface{}{
			"topic_filter": t.ToJSON(),
			"subscriber":   s.Address,
		},
	}

	_, err := s.socket.RequestRemoteService(&request)
	if err != nil {
		ch <- message.Fail(err.Error())
		return
	}

	s.timer = time.AfterFunc(time.Second*time.Duration(10), func() {
		ch <- message.Reply{Status: "fail", Message: "Server is not responding"}
	})

	go s.heartbeat(ch)

	return nil
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
	sub, err := remote.NewSub(s.socket.RemoteBroadcastUrl(), s.Address)
	if err != nil {
		return message.Fail("Failed to establish a connection with SDS Publisher: " + err.Error()), nil, nil
	}

	// now create a heartbeat timer
	ch := make(chan message.Broadcast)

	go s.loop(sub, ch)

	return message.Reply{Status: "OK", Message: "Successfully created a listener"}, ch, hb
}
