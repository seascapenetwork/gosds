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

func RemoteBlockRange(socket *remote.Socket, networkId string, address string, from uint64, to uint64) (uint64, []Transaction, []Log, error) {
	request := message.Request{
		Command: "block_get_range",
		Param: map[string]interface{}{
			"block_number_from": from,
			"block_number_to":   to,
			"to":                address,
			"network_id":        networkId,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, nil, nil, err
	}

	timestamp, err := message.GetUint64(params, "timestamp")
	if err != nil {
		return 0, nil, nil, err
	}

	raw_transactions, err := message.GetMapList(params, "transactions")
	if err != nil {
		return 0, nil, nil, err
	}

	raw_logs, err := message.GetMapList(params, "logs")
	if err != nil {
		return 0, nil, nil, err
	}

	transactions := make([]Transaction, len(raw_transactions))
	for i, raw := range raw_transactions {
		transactions[i] = ParseTransaction(raw)
	}

	logs := make([]Log, len(raw_logs))
	for i, raw := range raw_logs {
		logs[i] = ParseLog(raw)
	}

	return timestamp, transactions, logs, nil
}
