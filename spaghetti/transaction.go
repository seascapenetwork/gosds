/*Spaghetti transaction without method name and without clear input parameters*/
package spaghetti

import (
	"encoding/json"

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

func (b *Transaction) ToString() string {
	interfaces := b.ToJSON()
	byt, err := json.Marshal(interfaces)
	if err != nil {
		return ""
	}

	return string(byt)
}

func ParseTransaction(tx map[string]interface{}) Transaction {
	return Transaction{
		networkId:      tx["network_id"].(string),
		blockNumber:    int(tx["block_number"].(float64)),
		blockTimestamp: int(tx["block_timestamp"].(float64)),
		txid:           tx["txid"].(string),
		txIndex:        uint(tx["tx_index"].(float64)),
		txFrom:         tx["tx_from"].(string),
		txTo:           tx["tx_to"].(string),
		data:           tx["tx_data"].(string),
		value:          tx["tx_value"].(float64),
	}
}

func ParseTransactions(txs []interface{}) []Transaction {
	var transactions []Transaction
	for _, tx := range txs {
		if tx == nil {
			continue
		}
		transaction := ParseTransaction(tx.(map[string]interface{}))
		transactions = append(transactions, transaction)
	}
	return transactions
}

func RemoteTransactionDeployed(socket *remote.Socket, networkId string, txid string) (string, string, uint64, uint64, error) {
	// Send hello.
	request := message.Request{
		Command: "transaction_deployed_get",
		Param: map[string]interface{}{
			"network_id": networkId,
			"txid":       txid,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return "", "", 0, 0, err
	}

	address := params["address"].(string)
	deployer := params["deployer"].(string)
	blockNumber := uint64(params["block_number"].(float64))
	blockTimestamp := uint64(params["block_timestamp"].(float64))

	return address, deployer, blockNumber, blockTimestamp, nil
}
