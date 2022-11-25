package categorizer

import (
	"encoding/json"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
)

type Smartcontract struct {
	networkId         string
	address           string
	abiHash           string
	syncedBlockNumber int
	syncedTimestamp   int
}

func (b *Smartcontract) NetworkID() string {
	return b.networkId
}

func (b *Smartcontract) Key() string {
	return b.networkId + "." + b.address
}

func (b *Smartcontract) BlockNumber() int {
	return b.syncedBlockNumber
}

func (b *Smartcontract) Timestamp() int {
	return b.syncedTimestamp
}

func (b *Smartcontract) Address() string {
	return b.address
}

func (b *Smartcontract) AbiHash() string {
	return b.abiHash
}

func (b *Smartcontract) SetSyncing(n int, t int) {
	b.syncedBlockNumber = n
	b.syncedTimestamp = t
}

func New(networkId string, abiHash string, address string, syncedBlockNumber int, timestamp int) Smartcontract {
	return Smartcontract{
		networkId:         networkId,
		address:           address,
		abiHash:           abiHash,
		syncedBlockNumber: syncedBlockNumber,
		syncedTimestamp:   timestamp,
	}
}

func (b *Smartcontract) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = b.networkId
	i["address"] = b.address
	i["abi_hash"] = b.abiHash
	i["categorized_block_number"] = b.syncedBlockNumber
	i["categorized_block_timestamp"] = b.syncedTimestamp
	return i
}

func ParseJSON(blob map[string]interface{}) *Smartcontract {
	b := New(
		blob["network_id"].(string),
		blob["address"].(string),
		blob["abi_hash"].(string),
		int(blob["categorized_block_number"].(float64)),
		int(blob["categorized_block_timestamp"].(float64)),
	)
	return &b
}

func (b *Smartcontract) ToString() string {
	s := b.ToJSON()
	byt, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(byt)
}

func (b *Smartcontract) RemoteSmartcontractSet(socket *remote.Socket) error {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_set",
		Param:   b.ToJSON(),
	}

	_, err := socket.RequestRemoteService(&request)
	if err != nil {
		return err
	}

	return nil
}

func RemoteSmartcontract(socket *remote.Socket, networkId string, address string) (*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get",
		Param: map[string]interface{}{
			"network_id": networkId,
			"address":    address,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	b := ParseJSON(params["smartcontract"].(map[string]interface{}))
	return b, nil
}

func RemoteSmartcontracts(socket *remote.Socket) ([]*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get_all",
		Param:   map[string]interface{}{},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	returnedBlocks := params["smartcontracts"].([]interface{})
	smartcontracts := make([]*Smartcontract, len(returnedBlocks))

	for i, returnedBlock := range returnedBlocks {
		b := ParseJSON(returnedBlock.(map[string]interface{}))

		smartcontracts[i] = b
	}

	return smartcontracts, nil
}
