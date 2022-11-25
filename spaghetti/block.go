package spaghetti

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
)

// Returns the block minted time from SDS Spaghetti
func RemoteBlockMintedTime(socket *remote.Socket, networkId string, blockNumber uint64) (uint64, error) {
	// Send hello.
	request := message.Request{
		Command: "block_minted_time_get",
		Param: map[string]interface{}{
			"network_id":   networkId,
			"block_number": blockNumber,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, err
	}

	blockTimestamp := uint64(params["timestamp"].(float64))

	return blockTimestamp, nil
}
