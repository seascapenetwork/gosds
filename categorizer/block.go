package categorizer

import (
	"encoding/json"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
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

func (b *Block) RemoteSet(socket *remote.Socket) error {
	// Send hello.
	request := message.Request{
		Command: "categorizer_block_set",
		Param:   b.ToJSON(),
	}

	_, err := socket.RequestRemoteService(&request)
	if err != nil {
		return err
	}

	return nil
}

func RemoteBlock(socket *remote.Socket, networkId string, address string) (*Block, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get",
		Param: map[string]interface{}{
			"networkId": networkId,
			"address":   address,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	b := ParseJSON(params["block"].(map[string]interface{}))
	return b, nil
}

func RemoteBlocks(socket *remote.Socket) ([]*Block, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get_all",
		Param:   map[string]interface{}{},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	returnedBlocks := params["blocks"].([]interface{})
	blocks := make([]*Block, len(returnedBlocks))

	for i, returnedBlock := range returnedBlocks {
		b := ParseJSON(returnedBlock.(map[string]interface{}))

		blocks[i] = b
	}

	return blocks, nil
}
