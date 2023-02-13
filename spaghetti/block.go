package spaghetti

import (
	"fmt"
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"

	"github.com/blocklords/gosds/spaghetti"

	eth_types "github.com/ethereum/go-ethereum/core/types"
)

type Block struct {
	NetworkId      string
	BlockNumber    uint64
	BlockTimestamp uint64
	Transactions   []*spaghetti.Transaction
	Logs           []*spaghetti.Log
}

func (block *Block) SetTransactions(raw_transactions []*eth_types.Transaction) error {
	transactions := make([]*spaghetti.Transaction, len(raw_transactions))

	for txIndex, rawTx := range raw_transactions {
		tx, txErr := transaction.ParseTransaction(block.NetworkId, block.BlockNumber, uint(txIndex), rawTx)
		if txErr != nil {
			return fmt.Errorf("failed to set the block transactions. transaction parse error: %v", txErr)
		}

		transactions[txIndex] = tx
	}

	block.Transactions = transactions
	return nil
}

func (block *Block) SetLogs(raw_logs []eth_types.Log) error {
	var logs []*spaghetti.Log
	for _, rawLog := range raw_logs {
		if rawLog.Removed {
			continue
		}

		log, txErr := log.ParseLog(block.NetworkId, &rawLog)
		if txErr != nil {
			return txErr
		}

		logs = append(logs, log)
	}

	block.Logs = logs

	return nil
}

// Returns the earliest number in the cache for a given network id
func RemoteBlockEarliestNumber(socket *remote.Socket, network_id string) (uint64, error) {
	// Send hello.
	request := message.Request{
		Command: "block_get_earliest_cached_block_number",
		Parameters: map[string]interface{}{
			"network_id": network_id,
		},
	}

	parameters, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, err
	}

	return message.GetUint64(parameters, "block_number")
}

// Returns the block minted time from SDS Spaghetti
func RemoteBlockMintedTime(socket *remote.Socket, networkId string, blockNumber uint64) (uint64, error) {
	// Send hello.
	request := message.Request{
		Command: "block_get_timestamp",
		Parameters: map[string]interface{}{
			"network_id":   networkId,
			"block_number": blockNumber,
		},
	}

	parameters, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, err
	}

	return message.GetUint64(parameters, "block_timestamp")
}

func RemoteBlockRange(socket *remote.Socket, networkId string, address string, from uint64, to uint64) (uint64, []*Transaction, []*Log, error) {
	request := message.Request{
		Command: "block_get_range",
		Parameters: map[string]interface{}{
			"block_number_from": from,
			"block_number_to":   to,
			"to":                address,
			"network_id":        networkId,
		},
	}

	parameters, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, nil, nil, err
	}

	timestamp, err := message.GetUint64(parameters, "timestamp")
	if err != nil {
		return 0, nil, nil, err
	}

	raw_transactions, err := message.GetMapList(parameters, "transactions")
	if err != nil {
		return 0, nil, nil, err
	}

	raw_logs, err := message.GetMapList(parameters, "logs")
	if err != nil {
		return 0, nil, nil, err
	}

	transactions := make([]*Transaction, len(raw_transactions))
	for i, raw := range raw_transactions {
		tx, err := ParseTransaction(raw)
		if err != nil {
			return 0, nil, nil, err
		}
		transactions[i] = tx
	}

	logs := make([]*Log, len(raw_logs))
	for i, raw := range raw_logs {
		l, err := ParseLog(raw)
		if err != nil {
			return 0, nil, nil, err
		}
		logs[i] = l
	}

	return timestamp, transactions, logs, nil
}

// Returns the remote block information
// The address parameter is optional (make it a blank string)
// In that case SDS Spaghetti will return block with all transactions and logs.
func RemoteBlock(socket *remote.Socket, network_id string, block_number uint64, address string) (bool, uint64, []*Transaction, []*Log, error) {
	request := message.Request{
		Command: "block_get",
		Parameters: map[string]interface{}{
			"block_number": block_number,
			"network_id":   network_id,
			"to":           address,
		},
	}

	parameters, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, nil, nil, err
	}

	cached, err := message.GetBoolean(parameters, "cached")

	timestamp, err := message.GetUint64(parameters, "timestamp")
	if err != nil {
		return 0, nil, nil, err
	}

	raw_transactions, err := message.GetMapList(parameters, "transactions")
	if err != nil {
		return 0, nil, nil, err
	}

	raw_logs, err := message.GetMapList(parameters, "logs")
	if err != nil {
		return 0, nil, nil, err
	}

	transactions := make([]*Transaction, len(raw_transactions))
	for i, raw := range raw_transactions {
		tx, err := ParseTransaction(raw)
		if err != nil {
			return 0, nil, nil, err
		}
		transactions[i] = tx
	}

	logs := make([]*Log, len(raw_logs))
	for i, raw := range raw_logs {
		l, err := ParseLog(raw)
		if err != nil {
			return 0, nil, nil, err
		}
		logs[i] = l
	}

	return cached, timestamp, transactions, logs, nil
}
