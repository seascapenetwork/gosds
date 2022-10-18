package categorizer

type Transaction struct {
	ID          uint64
	NetworkId   string
	Address     string
	BlockNumber int
	Txid        string
	TxIndex     uint
	TxFrom      string
	Method      string
	Args        map[string]interface{}
	Value       float64
}

func (b *Transaction) Key() string {
	return b.NetworkId + "." + b.Address
}
