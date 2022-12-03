// The sdk/subscriber package is used to register for the smartcontracts
package subscriber

import (
	"fmt"
	"time"

	"github.com/blocklords/gosds/categorizer"
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/static"

	"github.com/blocklords/gosds/sdk/db"
)

type Subscriber struct {
	Address           string // Account address granted for reading
	timer             *time.Timer
	socket            *remote.Socket
	db                *db.KVM                    // it also keeps the topic filter
	smartcontractKeys []*static.SmartcontractKey // list of smartcontract keys

	BroadcastChan   chan message.Broadcast
	broadcastSocket *remote.Socket
}

func NewSubscriber(gatewaySocket *remote.Socket, db *db.KVM, address string) (*Subscriber, error) {
	subscriber := Subscriber{
		Address: address,
		socket:  gatewaySocket,
		db:      db,
	}

	return &subscriber, nil
}

func (subscriber *Subscriber) start_subscriber() error {
	// Run the Subscriber that is connected to the Broadcaster
	subscriber.broadcastSocket = remote.TcpSubscriberOrPanic(subscriber.socket.RemoteEnv())

	// Subscribing to the events, but we will not call the sub.ReceiveMessage
	// until we will not get the snapshot of the missing data.
	// ZMQ will queue the data until we will not call sub.ReceiveMessage.
	for _, key := range subscriber.smartcontractKeys {
		err := subscriber.broadcastSocket.SetSubscribeFilter(string(*key))
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

	if err := s.start_subscriber(); err != nil {
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
			tx, err := categorizer.ParseTransaction(rawTx)
			if err != nil {
				panic("failed to parse the transaction. the error: " + err.Error())
			} else {
				transactions[i] = tx
			}

			if uint64(transactions[i].BlockTimestamp) > latestBlockNumber {
				latestBlockNumber = uint64(transactions[i].BlockTimestamp)
			}
		}
		for i, rawLog := range rawLogs {
			log, err := categorizer.ParseLog(rawLog)
			if err != nil {
				panic("failed to parse the log. the error: " + err.Error())
			}
			logs[i] = log
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

	sock := remote.TcpPushSocketOrPanic(port)
	sock.SendMessage("")
	sock.Close()
}

func (s *Subscriber) runSink(port uint, smartcontractAmount int) {
	sock := remote.TcpPullSocketOrPanic(port)
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

// todo, change the heartbeat logic, expect to receive messages from the SDS Gateway
func (s *Subscriber) loop() {
	s.timer = time.AfterFunc(time.Second*time.Duration(10), func() {
		s.BroadcastChan <- message.NewBroadcast("", message.Reply{Status: "fail", Message: "Server is not responding"})
	})

	// todo
	// use remote/Subscribe
	// change heartbeat, upon expiration of the heartbeat start over.
	receive_channel := make(chan message.Reply)

	s.broadcastSocket.Subscribe(receive_channel, time.Second*30)

	for {
		reply := <-receive_channel

		if !reply.IsOK() {
			if reply.Message != "timeout" {
				//  Send results to sink
				s.BroadcastChan <- message.NewBroadcast("", reply)
				//  Exit, assume that the Client will restart it.
				// we might need to restart ourselves later.
				break
			} else {
				s.broadcastSocket.Close()

				if err := s.start_subscriber(); err != nil {
					fmt.Println("failed to start the subscriber")
					s.BroadcastChan <- message.NewBroadcast("error", message.Fail("failed to restart the subscriber: "+err.Error()))
					break
				}

				s.broadcastSocket.Subscribe(receive_channel, time.Second*30)

				break
			}
		}

		// we skip the duplicate messages that were fetched by the Snapshot
		networkId := reply.Params["network_id"].(string)
		address := reply.Params["address"].(string)
		blockTimestamp := uint64(reply.Params["block_timestamp"].(float64))
		key := static.CreateSmartcontractKey(networkId, address)
		latestBlockNumber := s.db.GetBlockTimestamp(&key)

		if latestBlockNumber > blockTimestamp {
			continue
		}

		s.db.SetBlockTimestamp(&key, blockTimestamp)
		s.BroadcastChan <- message.NewBroadcast("", reply)
	}
}
