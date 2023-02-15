// The network package is used to get the blockchain network information.
package static

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
)

type Network struct {
	Id       string
	Provider string
	Flag     int8 // With VM or Without VM
}

const (
	ALL        int8 = 0 // any blockchain
	WITH_VM    int8 = 1 // with EVM
	WITHOUT_VM int8 = 2 // without EVM, it's an L2
)

// Whether the given flag is valid Network Flag or not.
func IsValidFlag(flag int8) bool {
	return flag == WITH_VM || flag == WITHOUT_VM || flag == ALL
}

// parses JSON object into the Network Type
func ParseNetwork(raw map[string]interface{}) (*Network, error) {
	id, err := message.GetString(raw, "id")
	if err != nil {
		return nil, err
	}

	flag_64, err := message.GetUint64(raw, "flag")
	if err != nil {
		return nil, err
	}
	flag := int8(flag_64)
	if !IsValidFlag(flag) || flag == ALL {
		return nil, errors.New("invalid 'flag' from the parsed data")
	}

	provider, err := message.GetString(raw, "provider")
	if err != nil {
		return nil, err
	}

	return &Network{
		Id:       id,
		Provider: provider,
		Flag:     flag,
	}, nil
}

// parses list of JSON objects into the list of Networks
func ParseNetworks(raw_networks []map[string]interface{}) ([]*Network, error) {
	networks := make([]*Network, len(raw_networks))

	for i, raw := range raw_networks {
		network, err := ParseNetwork(raw)
		if err != nil {
			return nil, err
		}

		networks[i] = network
	}

	return networks, nil

func (n *Network) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"id":       n.Id,
		"provider": n.Provider,
		"flag":     n.Flag,
	}
}


			delete(supportedNetworks, networkId)
		}
	}

}

// Returns list of support network IDs from SDS Static
func GetNetworkIds(socket *remote.Socket, flag int8) ([]string, error) {
	if !IsValidFlag(flag) {
		return nil, errors.New("invalid 'flag' parameter")
	}
	request := message.Request{
		Command: "network_id_get_all",
		Parameters: map[string]interface{}{
			"flag": flag,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}
	return message.GetStringList(params, "network_ids")
}

// Returns list of support network IDs from SDS Static
func GetNetworks(socket *remote.Socket, flag int8) ([]*Network, error) {
	if !IsValidFlag(flag) {
		return nil, errors.New("invalid 'flag' parameter")
	}
	request := message.Request{
		Command: "network_get_all",
		Parameters: map[string]interface{}{
			"flag": flag,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}
	raw_networks, err := message.GetMapList(params, "networks")
	if err != nil {
		return nil, err
	}

	return ParseNetworks(raw_networks)
}

// Returns the Blockchain Network access provider
func GetNetwork(socket *remote.Socket, network_id string, flag int8) (*Network, error) {
	if !IsValidFlag(flag) {
		return nil, errors.New("invalid 'flag' parameter")
	}
	request := message.Request{
		Command: "network_get",
		Parameters: map[string]interface{}{
			"network_id": network_id,
			"flag":       flag,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}
	raw, err := message.GetMap(params, "network")
	if err != nil {
		return nil, err
	}

	return ParseNetwork(raw)
}
