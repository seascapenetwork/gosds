package categorizer

import "github.com/blocklords/gosds/spaghetti"

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
