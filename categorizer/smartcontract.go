package categorizer

import (
	"encoding/json"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
)

type Smartcontract struct {
	NetworkId                 string
	Address                   string
	CategorizedBlockNumber    uint64
	CategorizedBlockTimestamp uint64
}

func (s *Smartcontract) SmartcontractKeyString() string {
	return s.NetworkId + "." + s.Address
}

// Updates the categorized block parameter of the smartcontract.
// It means, this smartcontract 's' data was categorized till the given block numbers.
//
// The first is the block number, second is the block timestamp.
func (s *Smartcontract) SetBlockParameter(b uint64, t uint64) {
	s.CategorizedBlockNumber = b
	s.CategorizedBlockTimestamp = t
}

func (s *Smartcontract) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = s.NetworkId
	i["address"] = s.Address
	i["categorized_block_number"] = s.CategorizedBlockNumber
	i["categorized_block_timestamp"] = s.CategorizedBlockTimestamp
	return i
}

func ParseSmartcontract(blob map[string]interface{}) (*Smartcontract, error) {
	network_id, err := message.GetString(blob, "network_id")
	if err != nil {
		return nil, err
	}
	address, err := message.GetString(blob, "address")
	if err != nil {
		return nil, err
	}
	categorized_block_number, err := message.GetUint64(blob, "categorized_block_number")
	if err != nil {
		return nil, err
	}
	categorized_block_timestamp, err := message.GetUint64(blob, "categorized_block_timestamp")
	if err != nil {
		return nil, err
	}

	return &Smartcontract{
		NetworkId:                 network_id,
		Address:                   address,
		CategorizedBlockNumber:    categorized_block_number,
		CategorizedBlockTimestamp: categorized_block_timestamp,
	}, nil
}

// Returns a JSON representation of this smartcontract in a string format
func (b *Smartcontract) ToString() (string, error) {
	s := b.ToJSON()
	byt, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	return string(byt), nil
}

// Sends a command to the remote SDS Categorizer about regitration of this smartcontract.
func (b *Smartcontract) RemoteSet(socket *remote.Socket) error {
	// Send hello.
	request := message.Request{
		Command:    "smartcontract_set",
		Parameters: b.ToJSON(),
	}

	_, err := socket.RequestRemoteService(&request)
	if err != nil {
		return err
	}

	return nil
}

// Returns a smartcontract information from the remote SDS Categorizer.
func RemoteSmartcontract(socket *remote.Socket, network_id string, address string) (*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command: "smartcontract_get",
		Parameters: map[string]interface{}{
			"network_id": network_id,
			"address":    address,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	smartcontract, err := message.GetMap(params, "smartcontract")
	if err != nil {
		return nil, err
	}

	return ParseSmartcontract(smartcontract)
}

// Returns all smartcontracts from SDS Categorizer
func RemoteSmartcontracts(socket *remote.Socket) ([]*Smartcontract, error) {
	// Send hello.
	request := message.Request{
		Command:    "smartcontract_get_all",
		Parameters: map[string]interface{}{},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	raw_smartcontracts, err := message.GetMapList(params, "smartcontracts")
	if err != nil {
		return nil, err
	}

	smartcontracts := make([]*Smartcontract, len(raw_smartcontracts))
	for i, raw := range raw_smartcontracts {
		smartcontract, err := ParseSmartcontract(raw)
		if err != nil {
			return nil, err
		}

		smartcontracts[i] = smartcontract
	}

	return smartcontracts, nil
}
