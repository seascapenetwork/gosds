package categorizer

import (
	"encoding/json"
	"fmt"

	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

type Block struct {
	networkId         string
	address           string
	abiHash           string
	syncedBlockNumber int
	syncedTimestamp   int
}

func (b *Block) NetworkID() string {
	return b.networkId
}

func (b *Block) Key() string {
	return b.networkId + "." + b.address
}

func (b *Block) BlockNumber() int {
	return b.syncedBlockNumber
}

func (b *Block) Timestamp() int {
	return b.syncedTimestamp
}

func (b *Block) Address() string {
	return b.address
}

func (b *Block) AbiHash() string {
	return b.abiHash
}

func (b *Block) SetSyncing(n int, t int) {
	b.syncedBlockNumber = n
	b.syncedTimestamp = t
}

func New(networkId string, abiHash string, address string, syncedBlockNumber int, timestamp int) Block {
	return Block{
		networkId:         networkId,
		address:           address,
		abiHash:           abiHash,
		syncedBlockNumber: syncedBlockNumber,
		syncedTimestamp:   timestamp,
	}
}

func (b *Block) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = b.networkId
	i["address"] = b.address
	i["abi_hash"] = b.abiHash
	i["synced_block_number"] = b.syncedBlockNumber
	i["timestamp"] = b.syncedTimestamp
	return i
}

func ParseJSON(blob map[string]interface{}) *Block {
	b := New(
		blob["network_id"].(string),
		blob["address"].(string),
		blob["abi_hash"].(string),
		int(blob["synced_block_number"].(float64)),
		int(blob["timestamp"].(float64)),
	)
	return &b
}

func (b *Block) ToString() string {
	s := b.ToJSON()
	byt, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(byt)
}

func (b *Block) Save(socket *zmq.Socket) error {
	// Send hello.
	request := message.Request{
		Command: "categorizer_block_set",
		Param:   b.ToJSON(),
	}

	if _, err := socket.SendMessage(request.ToString()); err != nil {
		return fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		return fmt.Errorf("receiving: %w", err)
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		return fmt.Errorf("sds categorizer reply for setting new block: %w", err)
	}
	if !reply.IsOK() {
		return fmt.Errorf("sds categorizer reply for setting new block not ok: %s", reply.Message)
	}

	return nil
}

func RemoteGet(socket *zmq.Socket, networkId string, address string) (*Block, error) {
	// Send hello.
	request := message.Request{
		Command: "get",
		Param: map[string]interface{}{
			"network_id": networkId,
			"address":    address,
		},
	}
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		return nil, fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		return nil, fmt.Errorf("receiving from SDS Categorizer for 'get' command: %w", err)
	}

	fmt.Println(r)
	reply, err := message.ParseReply(r)
	if err != nil {
		return nil, fmt.Errorf("categorizer block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		return nil, fmt.Errorf("categorizer block reply status is not ok: %s", reply.Message)
	}

	b := ParseJSON(reply.Params["block"].(map[string]interface{}))
	return b, nil
}

func RemoteGetAll(socket *zmq.Socket) ([]*Block, error) {
	// Send hello.
	request := message.Request{
		Command: "get_all",
		Param:   map[string]interface{}{},
	}
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to get all the blocks from SDS-Categorizer", err.Error())
		return nil, fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller")
		return nil, fmt.Errorf("receiving: %w", err)
	}

	fmt.Println(r)
	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse abi reply")
		return nil, fmt.Errorf("spaghetti block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure")
		return nil, fmt.Errorf("spaghetti block reply status is not ok: %s", reply.Message)
	}

	returnedBlocks := reply.Params["blocks"].([]interface{})
	blocks := make([]*Block, len(returnedBlocks))
	for i, returnedBlock := range returnedBlocks {
		b := ParseJSON(returnedBlock.(map[string]interface{}))

		blocks[i] = b
	}

	return blocks, nil
}
