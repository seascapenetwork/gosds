package categorizer

import (
	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/spaghetti"
)

type Transaction struct {
	ID             string
	NetworkId      string
	Address        string
	BlockNumber    int
	BlockTimestamp int
	Txid           string
	TxIndex        uint
	TxFrom         string
	Method         string
	Args           map[string]interface{}
	Value          float64
}

func TransactionKey(networkId string, txId string) string {
	return networkId + "." + txId
}

func (b *Transaction) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = b.NetworkId
	i["address"] = b.Address
	i["block_number"] = b.BlockNumber
	i["block_timestamp"] = b.BlockTimestamp
	i["txid"] = b.Txid
	i["tx_index"] = b.TxIndex
	i["tx_from"] = b.TxFrom
	i["method"] = b.Method
	i["arguments"] = b.Args
	i["value"] = b.Value
	return i
}

func ParseTransactionFromJson(blob map[string]interface{}) *Transaction {
	return &Transaction{
		NetworkId:      blob["network_id"].(string),
		Address:        blob["address"].(string),
		BlockNumber:    int(blob["block_number"].(float64)),
		BlockTimestamp: int(blob["block_timestamp"].(float64)),
		Txid:           blob["txid"].(string),
		TxIndex:        blob["tx_index"].(uint),
		TxFrom:         blob["tx_from"].(string),
		Method:         blob["method"].(string),
		Args:           blob["arguments"].(map[string]interface{}),
		Value:          blob["value"].(float64),
	}
}

func ParseTransaction(tx spaghetti.Transaction, method string, inputs map[string]interface{}, c *Block, blockNumber int, blockTimestamp int) Transaction {
	return Transaction{
		NetworkId:      c.NetworkID(),
		Address:        c.Address(),
		BlockNumber:    blockNumber,
		BlockTimestamp: blockTimestamp,
		Txid:           tx.TxId(),
		TxIndex:        tx.TxIndex(),
		TxFrom:         tx.TxFrom(),
		Method:         method,
		Args:           inputs,
		Value:          tx.Value(),
	}
}

func RemoteTransactionAmount(socket *remote.Socket, blockTimestampFrom int, blockTimestampTo int, smartcontractKeys []string) (int, error) {
	request := message.Request{
		Command: "transaction_amount",
		Param: map[string]interface{}{
			"block_timestamp_from": blockTimestampFrom,
			"block_timestamp_to":   blockTimestampTo,
			"smartcontract_keys":   smartcontractKeys,
		},
	}
	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return 0, err
	}

	txAmount := int(params["transaction_amount"].(float64))

	return txAmount, nil
}

func RemoteTransactions(socket *remote.Socket, blockTimestampFrom int, blockTimestampTo int, smartcontractKeys []string, page int, limit uint) ([]*Transaction, error) {
	request := message.Request{
		Command: "transaction_get_all",
		Param: map[string]interface{}{
			"block_timestamp_from": blockTimestampFrom,
			"block_timestamp_to":   blockTimestampTo,
			"smartcontract_keys":   smartcontractKeys,
			"page":                 page,
			"limit":                limit,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	raws := params["transactions"].([]interface{})
	transactions := make([]*Transaction, len(raws))
	for i, raw := range raws {
		transactions[i] = ParseTransactionFromJson(raw.(map[string]interface{}))
	}

	return transactions, nil
}
