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

// the main function that starts the broadcasting.
// It first calls the smartcontract_filters. and cacshes them out.
// if there is an error, it will return them either in the Heartbeat channel
func (s *Subscriber) Listen(t *topic.TopicFilter) (message.Reply, chan message.Broadcast, chan message.Reply) {
	hb := make(chan message.Reply)

	err := s.initiate(hb)
	if err != nil {
		return message.Fail("subscribe initiation error: " + err.Error()), nil, nil
	}

	// Run the listener
	sub, err := remote.NewSub(s.socket.RemoteBroadcastUrl(), s.Address)
	if err != nil {
		return message.Fail("failed to establish a connection with SDS Gateway: " + err.Error()), nil, nil
	}
	// Subscribing to the events, but we will not call the sub.ReceiveMessage
	// until we will not get the snapshot of the missing data.
	for _, key := range s.smartcontractKeys {
		err := sub.SetSubscribe(string(*key))
		if err != nil {
			return message.Fail("failed to subscribe to the smartcontract: " + err.Error()), nil, nil
		}
	}

	// now create a heartbeat timer
	ch := make(chan message.Broadcast)

	go s.loop(sub, ch)

	return message.Reply{Status: "OK", Message: "Successfully created a listener"}, ch, hb
}

// The algorithm
// List of the smartcontracts by smartcontract filter
func (s *Subscriber) initiate(ch chan message.Reply) error {
	// preparing the subscriber so that we catch the first message if it was send
	// by publisher.

	smartcontracts, topicStrings, err := static.RemoteSmartcontracts(s.socket, s.db.TopicFilter())
	if err != nil {
		ch <- message.Fail(err.Error())
		return err
	}

	// set the smartcontract keys
	for i, sm := range smartcontracts {
		key := sm.KeyString()

		// cache the smartcontract block timestamp
		// block timestamp is used to subscribe for the events
		blockTimestamp := s.db.GetBlockTimestamp(&key)
		if blockTimestamp == 0 {
			blockTimestamp = uint64(sm.PreDeployBlockTimestamp)
			err := s.db.SetBlockTimestamp(&key, blockTimestamp)
			if err != nil {
				ch <- message.Fail(err.Error())
				return err
			}
		}

		// cache the topic string
		topicString := topicStrings[i]
		err := s.db.SetTopicString(&key, topicString)
		if err != nil {
			ch <- message.Fail(err.Error())
			return err
		}

		// finally track the smartcontract
		s.smartcontractKeys = append(s.smartcontractKeys, &key)
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
