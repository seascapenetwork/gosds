// The sdk/subscriber package is used to register for the smartcontracts
package subscriber

import (
	"fmt"
	"time"

	"github.com/blocklords/gosds/categorizer"
	"github.com/blocklords/gosds/generic_type"
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/static"

	"github.com/blocklords/gosds/sdk/db"
)

type Subscriber struct {
	Address           string // Account address granted for reading
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
	err := s.load_smartcontracts()
	if err != nil {
		return err
	}

	if err := s.start_subscriber(); err != nil {
		return err
	}


	// now create a broadcaster channel to send back to the developer the messages
	s.BroadcastChan = make(chan message.Broadcast)


	go s.snapshot()
	return nil
}

// Returns the latest updated block timestamp in the cache
func (s *Subscriber) recent_block_timestamp() uint64 {
	var recent_block_timestamp uint64 = 0
	for _, key := range s.smartcontractKeys {
		block_timestamp := s.db.GetBlockTimestamp(*key)
		fmt.Println("recent block timestamp: ", *key, block_timestamp)
		if block_timestamp > recent_block_timestamp {
			recent_block_timestamp = block_timestamp
		}
	}

	return recent_block_timestamp
}

func (s *Subscriber) snapshot() {
	limit := uint64(500)
	page := uint64(1)
	blockTimestampFrom := s.recent_block_timestamp()
	blockTimestampTo := uint64(0)

	for {
		request := message.Request{
			Command: "snapshot_get",
			Param: map[string]interface{}{
				"smartcontract_keys":   generic_type.ToStringList(s.smartcontractKeys),
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

		raw_transactions, err := message.GetMapList(replyParams, "transactions")
		if err != nil {
			panic(err)
		}
		raw_logs, err := message.GetMapList(replyParams, "logs")
		if err != nil {
			panic(err)
		}
		timestamp, err := message.GetUint64(replyParams, "block_timestamp")
		if err != nil {
			panic(err)
		}

		// we fetch until all is not received
		if len(raw_transactions) == 0 {
			break
		}

		transactions := make([]*categorizer.Transaction, len(raw_transactions))
		logs := make([]*categorizer.Log, len(raw_logs))

		// Saving the latest block number in the cache
		// along the parsing raw data into SDS data type
		for i, rawTx := range raw_transactions {
			tx, err := categorizer.ParseTransaction(rawTx)
			if err != nil {
				panic("failed to parse the transaction. the error: " + err.Error())
			} else {
				transactions[i] = tx
			}

			key := static.CreateSmartcontractKey(tx.NetworkId, tx.Address)
			cached_block_timestamp := s.db.GetBlockTimestamp(key)
			if tx.BlockTimestamp > cached_block_timestamp {
				err = s.db.SetBlockTimestamp(key, tx.BlockTimestamp)
				if err != nil {
					panic(err)
				}
			}
		}
		for i, rawLog := range raw_logs {
			log, err := categorizer.ParseLog(rawLog)
			if err != nil {
				panic("failed to parse the log. the error: " + err.Error())
			}
			logs[i] = log
		}

		reply := message.Reply{
			Status:  "OK",
			Message: "",
			Params: map[string]interface{}{
				"transactions":    transactions,
				"logs":            logs,
				"block_timestamp": timestamp,
			},
		}
		s.BroadcastChan <- message.NewBroadcast("", reply)

		if blockTimestampTo == 0 {
			blockTimestampTo = timestamp
		}
		page++

	}

	s.loop()
}

// Get the list of the smartcontracts by smartcontract filter from SDS Categorizer via SDS Gateway
// Then cache them out and list in the Subscriber data structure
func (s *Subscriber) load_smartcontracts() error {
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
		blockTimestamp := s.db.GetBlockTimestamp(key)

		if blockTimestamp == 0 {
			blockTimestamp = uint64(sm.PreDeployBlockTimestamp)
			err := s.db.SetBlockTimestamp(key, blockTimestamp)
			if err != nil {
				return err
			}
		}

		// cache the topic string
		topicString := topicStrings[i]
		err = s.db.SetTopicString(key, topicString)
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
	receive_channel := make(chan message.Reply)
	exit_channel := make(chan int)

	time_out := time.Duration(time.Second * 30)

	go s.broadcastSocket.Subscribe(receive_channel, exit_channel, time_out)

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
				err := s.broadcastSocket.Close()
				if err != nil {
					panic(err)
				}

				if err := s.start_subscriber(); err != nil {
					fmt.Println("failed to start the subscriber")
					s.BroadcastChan <- message.NewBroadcast("error", message.Fail("failed to restart the subscriber: "+err.Error()))
					break
				}

				go s.broadcastSocket.Subscribe(receive_channel, exit_channel, time_out)
			}
			continue
		}

		// we skip the duplicate messages that were fetched by the Snapshot
		networkId, err := message.GetString(reply.Params, "network_id")
		if err != nil {
			fmt.Println("failed to receive the 'network_id' from the SDS Gateway Broadcast Proxy")
			fmt.Println("skip it. which we should not actually.")
			continue
		}
		address, err := message.GetString(reply.Params, "address")
		if err != nil {
			fmt.Println("failed to receive the 'address' from the SDS Gateway Broadcast Proxy")
			fmt.Println("skip it. which we should not actually.")
			continue
		}
		block_timestamp, err := message.GetUint64(reply.Params, "block_timestamp")
		if err != nil {
			fmt.Println("failed to receive the network_id from the SDS Gateway Broadcast Proxy")
			fmt.Println("skip it. which we should not actually.")
			continue
		}
		key := static.CreateSmartcontractKey(networkId, address)
		latestBlockNumber := s.db.GetBlockTimestamp(key)

		if latestBlockNumber > block_timestamp {
			continue
		}

		err = s.db.SetBlockTimestamp(key, block_timestamp)
		if err != nil {
			fmt.Println("failed to cache the block timestamp")
			fmt.Println("skip it. which we should not actually.")
			continue
		}
		s.BroadcastChan <- message.NewBroadcast("", reply)
	}
}
