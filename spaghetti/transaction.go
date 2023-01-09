/*Spaghetti transaction without method name and without clear input parameters*/
package spaghetti

import (
	"encoding/json"
	"errors"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
)

type Transaction struct {
	NetworkId      string
	BlockNumber    uint64
	BlockTimestamp uint64
	Txid           string // txId column
	TxFrom         string
	TxTo           string
	TxIndex        uint
	Data           string  // text Data type
	Value          float64 // Value attached with transaction
}

// JSON representation of the spaghetti.Transaction
func (b *Transaction) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"network_id":      b.NetworkId,
		"block_number":    b.BlockNumber,
		"block_timestamp": b.BlockTimestamp,
		"Txid":            b.Txid,
		"tx_from":         b.TxFrom,
		"tx_to":           b.TxTo,
		"tx_index":        b.TxIndex,
		"tx_Data":         b.Data,
		"tx_Value":        b.Value,
	}
}

// JSON string representation of the spaghetti.Transaction
func (b *Transaction) ToString() string {
	interfaces := b.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return ""
	}

	return string(byt)
}

// Parse the JSON into spaghetti.Transation
func ParseTransaction(parameters map[string]interface{}) (*Transaction, error) {
	network_id, err := message.GetString(parameters, "network_id")
	if err != nil {
		return nil, err
	}
	block_number, err := message.GetUint64(parameters, "block_number")
	if err != nil {
		return nil, err
	}
	block_timestamp, err := message.GetUint64(parameters, "block_timestamp")
	if err != nil {
		return nil, err
	}
	Txid, err := message.GetString(parameters, "Txid")
	if err != nil {
		return nil, err
	}
	tx_index, err := message.GetUint64(parameters, "tx_index")
	if err != nil {
		return nil, err
	}
	tx_from, err := message.GetString(parameters, "tx_from")
	if err != nil {
		return nil, err
	}
	tx_to, err := message.GetString(parameters, "tx_to")
	if err != nil {
		return nil, err
	}
	tx_Data, err := message.GetString(parameters, "tx_Data")
	if err != nil {
		return nil, err
	}
	Value, err := message.GetFloat64(parameters, "tx_Value")
	if err != nil {
		return nil, err
	}

	return &Transaction{
		NetworkId:      network_id,
		BlockNumber:    block_number,
		BlockTimestamp: block_timestamp,
		Txid:           Txid,
		TxIndex:        uint(tx_index),
		TxFrom:         tx_from,
		TxTo:           tx_to,
		Data:           tx_Data,
		Value:          Value,
	}, nil
}

func ParseTransactions(txs []interface{}) ([]*Transaction, error) {
	var transactions []*Transaction = make([]*Transaction, len(txs))
	for i, raw := range txs {
		if raw == nil {
			continue
		}
		map_log, ok := raw.(map[string]interface{})
		if !ok {
			return nil, errors.New("transaction is not a map")
		}
		transaction, err := ParseTransaction(map_log)
		if err != nil {
			return nil, err
		}
		transactions[i] = transaction
	}
	return transactions, nil
}

// Sends the command to the remote SDS Spaghetti to get the smartcontract deploy metaData by
// its transaction id
func RemoteTransactionDeployed(socket *remote.Socket, network_id string, Txid string) (string, string, uint64, uint64, error) {
	// Send hello.
	request := message.Request{
		Command: "transaction_deployed_get",
		Parameters: map[string]interface{}{
			"network_id": network_id,
			"txid":       Txid,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return "", "", 0, 0, err
	}

	address, err := message.GetString(params, "address")
	if err != nil {
		return "", "", 0, 0, err
	}
	deployer, err := message.GetString(params, "deployer")
	if err != nil {
		return "", "", 0, 0, err
	}
	block_number, err := message.GetUint64(params, "block_number")
	if err != nil {
		return "", "", 0, 0, err
	}
	block_timestamp, err := message.GetUint64(params, "block_timestamp")
	if err != nil {
		return "", "", 0, 0, err
	}

	return address, deployer, block_number, block_timestamp, nil
}
