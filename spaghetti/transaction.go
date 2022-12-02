/*Spaghetti transaction without method name and without clear input parameters*/
package spaghetti

import (
	"encoding/json"
	"errors"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
)

type Transaction struct {
	networkId      string
	blockNumber    int
	blockTimestamp int
	txid           string // txId column
	txFrom         string
	txTo           string
	txIndex        uint
	data           string  // text data type
	value          float64 // value attached with transaction
}

func (b *Transaction) NetworkID() string {
	return b.networkId
}

func (b *Transaction) BlockNumber() int {
	return b.blockNumber
}

func (b *Transaction) TxId() string {
	return b.txid
}

func (b *Transaction) TxFrom() string {
	return b.txFrom
}

func (b *Transaction) TxTo() string {
	return b.txTo
}

func (b *Transaction) TxIndex() uint {
	return b.txIndex
}

func (b *Transaction) Data() string {
	return b.data
}

func (b *Transaction) Value() float64 {
	return b.value
}

func (b *Transaction) Timestamp() int {
	return b.blockTimestamp
}

// JSON representation of the spaghetti.Transaction
func (b *Transaction) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"network_id":      b.networkId,
		"block_number":    b.blockNumber,
		"block_timestamp": b.blockTimestamp,
		"txid":            b.txid,
		"tx_from":         b.txFrom,
		"tx_to":           b.txTo,
		"tx_index":        b.txIndex,
		"tx_data":         b.data,
		"tx_value":        b.value,
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
	txid, err := message.GetString(parameters, "txid")
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
	tx_data, err := message.GetString(parameters, "tx_data")
	if err != nil {
		return nil, err
	}
	value, err := message.GetFloat64(parameters, "tx_value")
	if err != nil {
		return nil, err
	}

	return &Transaction{
		networkId:      network_id,
		blockNumber:    int(block_number),
		blockTimestamp: int(block_timestamp),
		txid:           txid,
		txIndex:        uint(tx_index),
		txFrom:         tx_from,
		txTo:           tx_to,
		data:           tx_data,
		value:          value,
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

// Sends the command to the remote SDS Spaghetti to get the smartcontract deploy metadata by
// its transaction id
func RemoteTransactionDeployed(socket *remote.Socket, network_id string, txid string) (string, string, uint64, uint64, error) {
	// Send hello.
	request := message.Request{
		Command: "transaction_deployed_get",
		Param: map[string]interface{}{
			"network_id": network_id,
			"txid":       txid,
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
