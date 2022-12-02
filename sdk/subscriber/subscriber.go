// The sdk/subscriber package is used to register for the smartcontracts
package subscriber

import (
	"fmt"
	"time"

	"github.com/blocklords/gosds/categorizer"
	"github.com/blocklords/gosds/message"
	sds_remote "github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/static"

	"github.com/blocklords/gosds/sdk/db"
	"github.com/blocklords/gosds/sdk/remote"

	zmq "github.com/pebbe/zmq4"
)

type Subscriber struct {
<<<<<<< HEAD
	Host                 string // SDS Gateway host
	Sub                  string // SDS Publisher host
	Address              string // Account address granted for reading
	timer                *time.Timer
	GatewaySocketContext *zmq.Socket
	SubSocketContext     *zmq.Socket
	StateTimestamp       int
	connectClosed        bool
=======
	Address           string // Account address granted for reading
	timer             *time.Timer
	socket            *sds_remote.Socket
	db                *db.KVM                    // it also keeps the topic filter
	smartcontractKeys []*static.SmartcontractKey // list of smartcontract keys

	BroadcastChan   chan message.Broadcast
	broadcastSocket *zmq.Socket
>>>>>>> main
}

func NewSubscriber(gatewaySocket *sds_remote.Socket, db *db.KVM, address string) (*Subscriber, error) {
	subscriber := Subscriber{
		Address: address,
		socket:  gatewaySocket,
		db:      db,
	}

	return &subscriber, nil
}

func (subscriber *Subscriber) startSubscriber() error {
	var err error

	// Run the Subscriber that is connected to the Broadcaster
	subscriber.broadcastSocket, err = remote.NewSub(subscriber.socket.RemoteBroadcastUrl(), subscriber.Address)
	if err != nil {
		return fmt.Errorf("failed to establish a connection with SDS Gateway: " + err.Error())
	}

	// Subscribing to the events, but we will not call the sub.ReceiveMessage
	// until we will not get the snapshot of the missing data.
	// ZMQ will queue the data until we will not call sub.ReceiveMessage.
	for _, key := range subscriber.smartcontractKeys {
		err := subscriber.broadcastSocket.SetSubscribe(string(*key))
		if err != nil {
			subscriber.broadcastSocket.Close()
			return fmt.Errorf("failed to subscribe to the smartcontract: " + err.Error())
		}
	}

	return nil
}

// the main function that starts the broadcasting.
// It first calls the smartcontract_filters. and cacshes them out.
// if there is an error, it will return them either in the Heartbeat channel
func (s *Subscriber) Start() error {
	s.smartcontractKeys = make([]*static.SmartcontractKey, 0)
	err := s.loadSmartcontracts()
	if err != nil {
		return err
	}

	if err := s.startSubscriber(); err != nil {
		return err
	}

	port, err := s.getSinkPort()
	if err != nil {
		s.broadcastSocket.Close()
		return fmt.Errorf("failed to create a port for Snapshots")
	}
	// run the think that waits for the snapshots
	// then it will start to receive messages from subscriber
	go s.runSink(port, len(s.smartcontractKeys))

	// now create a broadcaster channel to send back to the developer the messages
	s.BroadcastChan = make(chan message.Broadcast)

	// finally take the snapshots
	for i := range s.smartcontractKeys {
		go s.snapshot(i, port)
	}

	return nil
}

func (s *Subscriber) getSinkPort() (uint, error) {
	port, err := s.socket.RemoteBroadcastPort()
	if err != nil {
		return 0, err
	}

	return port + 1, nil
}

func (s *Subscriber) snapshot(i int, port uint) {
	key := s.smartcontractKeys[i]
	limit := uint64(500)
	page := uint64(1)
	blockTimestampFrom := s.db.GetBlockTimestamp(key)
	blockTimestampTo := uint64(0)

	for {
		request := message.Request{
			Command: "snapshot_get",
			Param: map[string]interface{}{
				"smartcontract_key":    key,
				"block_timestamp_from": blockTimestampFrom,
				"block_timestamp_to":   blockTimestampTo,
				"page":                 page,
				"limit":                limit,
			},
		}

		replyParams, err := s.socket.RequestRemoteService(&request)
		if err != nil {
			panic(err)
		}

		rawTransactions := replyParams["transactions"].([]map[string]interface{})
		rawLogs := replyParams["logs"].([]map[string]interface{})
		timestamp := uint64(replyParams["block_timestamp"].(float64))

		// we fetch until all is not received
		if len(rawTransactions) == 0 {
			break
		}

		transactions := make([]*categorizer.Transaction, len(rawTransactions))
		logs := make([]*categorizer.Log, len(rawLogs))

		latestBlockNumber := uint64(0)
		for i, rawTx := range rawTransactions {
			transactions[i] = categorizer.ParseTransactionFromJson(rawTx)

			if uint64(transactions[i].BlockTimestamp) > latestBlockNumber {
				latestBlockNumber = uint64(transactions[i].BlockTimestamp)
			}
		}
		for i, rawLog := range rawLogs {
			logs[i] = categorizer.ParseLog(rawLog)
		}

		err = s.db.SetBlockTimestamp(key, latestBlockNumber)
		if err != nil {
			panic(err)
		}

		reply := message.Reply{Status: "OK", Message: "", Params: map[string]interface{}{
			"transactions":    transactions,
			"logs":            logs,
			"block_timestamp": timestamp,
		}}
		s.BroadcastChan <- message.NewBroadcast(string(*key), reply)

		if blockTimestampTo == 0 {
			blockTimestampTo = timestamp
		}
		page++
	}

	sock := sds_remote.TcpPushSocketOrPanic(port)
	sock.SendMessage("")
	sock.Close()
}

func (s *Subscriber) runSink(port uint, smartcontractAmount int) {
	sock := sds_remote.TcpPullSocketOrPanic(port)
	defer sock.Close()

	for {
		_, err := sock.RecvMessage(0)
		if err != nil {
			fmt.Println("failed to receive a message: ", err.Error())
			continue
		}

		smartcontractAmount--

		if smartcontractAmount == 0 {
			break
		}
	}

	go s.loop()
}

// The algorithm
// Get the list of the smartcontracts by smartcontract filter from SDS Categorizer via SDS Gateway
// Then cache them out and list in the Subscriber data structure
func (s *Subscriber) loadSmartcontracts() error {
	// preparing the subscriber so that we catch the first message if it was send
	// by publisher.

	smartcontracts, topicStrings, err := static.RemoteSmartcontracts(s.socket, s.db.TopicFilter())
	if err != nil {
		return err
	}

<<<<<<< HEAD
	subscribed := remote.ReqReply(s.Host, msg)
	if !subscribed.IsOK() {
		fmt.Println(subscribed)
		ch <- subscribed
		return
	}

	s.StateTimestamp = int(subscribed.Params["block_timestamp"].(float64))

	s.timer = time.AfterFunc(time.Second*time.Duration(10), func() {
		ch <- message.Reply{Status: "fail", Message: "Server is not responding"}
	})

	go s.heartbeat(ch)
}

func (s *Subscriber) heartbeat(ch chan message.Reply) {
	msg := message.Request{
		Command: "heartbeat",
		Param: map[string]interface{}{
			"subscriber": s.Address,
		},
	}

=======
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
				return err
			}
		}

		// cache the topic string
		topicString := topicStrings[i]
		err := s.db.SetTopicString(&key, topicString)
		if err != nil {
			return err
		}

		// finally track the smartcontract
		s.smartcontractKeys = append(s.smartcontractKeys, &key)
	}

	return nil
}

func (s *Subscriber) heartbeat() {
>>>>>>> main
	for {
		if s.connectClosed == true {
			fmt.Println("!!! socket lost, break loop for heartbeat")
			break
		}

		s.timer.Reset(time.Second * time.Duration(10))

		heartbeatReply := remote.ReqReplyBySockContext(s.GatewaySocketContext, msg)
		if !heartbeatReply.IsOK() {
			s.BroadcastChan <- message.NewBroadcast("", heartbeatReply)
			break
		}

		stateTimestamp := int(heartbeatReply.Params["stateTimestamp"].(float64))
		if stateTimestamp != s.StateTimestamp {
			exception := message.Reply{
				Status:  "NOTOK",
				Message: "Timestamp is not same as Server, you should re-connect",
			}
			ch <- exception
		}

		//heartbeatReply := Heartbeat(s)
		//if !heartbeatReply.IsOK() {
		//	ch <- heartbeatReply
		//	break
		//}

		time.Sleep(time.Second * time.Duration(2))
	}
}

// todo, change the heartbeat logic, expect to receive messages from the SDS Gateway
func (s *Subscriber) loop() {
	s.timer = time.AfterFunc(time.Second*time.Duration(10), func() {
		s.BroadcastChan <- message.NewBroadcast("", message.Reply{Status: "fail", Message: "Server is not responding"})
	})

	go s.heartbeat()

	for {
<<<<<<< HEAD
		fmt.Println(s.connectClosed)
		if s.connectClosed == true {
			fmt.Println("!!! socket lost, break loop for receive message loop")
			break
		}

		msg_raw, err := sub.RecvMessage(0)
=======
		msg_raw, err := s.broadcastSocket.RecvMessage(0)
>>>>>>> main
		if err != nil {
			s.BroadcastChan <- message.NewBroadcast("", message.Reply{Status: "fail", Message: "receive error: " + err.Error()})
			break
		}
		// empty messages are skipped
		if len(msg_raw) == 0 {
			continue
		}

		b, err := message.ParseBroadcast(msg_raw)
		if err != nil {
			s.BroadcastChan <- message.NewBroadcast(s.Address, message.Reply{Status: "fail", Message: "Error when parsing message " + err.Error()})
			break
		}

		if !b.IsOK() {
			//  Send results to sink
			s.BroadcastChan <- b
			//  Exit, assume that the Client will restart it.
			// we might need to restart ourselves later.
			break
		}

		// we skip the duplicate messages that were fetched by the Snapshot
		networkId := b.Reply().Params["network_id"].(string)
		address := b.Reply().Params["address"].(string)
		blockTimestamp := uint64(b.Reply().Params["block_timestamp"].(float64))
		key := static.CreateSmartcontractKey(networkId, address)
		latestBlockNumber := s.db.GetBlockTimestamp(&key)

		if latestBlockNumber > blockTimestamp {
			continue
		}

		s.db.SetBlockTimestamp(&key, blockTimestamp)
		s.BroadcastChan <- b
	}
}
<<<<<<< HEAD

func (s *Subscriber) Listen(t *topic.TopicFilter) (message.Reply, chan message.Broadcast, chan message.Reply) {
	fmt.Println("Listen...")
	hb := make(chan message.Reply)

	go s.subscribe(t, hb)

	fmt.Println("connect to gateway: ", s.Host)
	gate, err := remote.NewSocket(s.Host)
	if err != nil {
		return message.Fail("Failed to establish a connection with SDS Gateway: " + err.Error()), nil, nil
	}
	s.GatewaySocketContext = gate

	// Run the listener
	fmt.Println("connect to publisher: ", s.Sub)
	sub, err := remote.NewSub(s.Sub, s.Address)
	if err != nil {
		return message.Fail("Failed to establish a connection with SDS Publisher: " + err.Error()), nil, nil
	}
	s.SubSocketContext = sub

	s.connectClosed = false

	// now create a heartbeat timer
	ch := make(chan message.Broadcast)

	go s.loop(sub, ch)

	return message.Reply{Status: "OK", Message: "Successfully created a listener"}, ch, hb
}

func (s *Subscriber) Close() {
	//if s.SubSocketContext != nil {
	//	fmt.Println("@@@ close sub socket @@@")
	//	s.SubSocketContext.Close()
	//}
	//if s.GatewaySocketContext != nil {
	//	fmt.Println("@@@ close gateway socket @@@")
	//	s.GatewaySocketContext.Close()
	//}
	s.connectClosed = true
}
=======
>>>>>>> main
