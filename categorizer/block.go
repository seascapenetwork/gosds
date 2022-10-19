package categorizer

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/blocklords/gosds/message"

	zmq "github.com/pebbe/zmq4"
)

type Block struct {
	networkId         string
	address           string
	abiHash           string
	syncedBlockNumber int
}

func (b *Block) NetworkID() string {
	return b.networkId
}

func (b *Block) Key() string {
	return b.networkId + "." + b.address
}

func (b *Block) BlockNumber() int {
	return b.syncedBlockNumber
}

func (b *Block) Address() string {
	return b.address
}

func (b *Block) AbiHash() string {
	return b.abiHash
}

func (b *Block) SetBlockNumber(stmt *sql.Stmt, n int) bool {
	b.syncedBlockNumber = n

	if _, err := stmt.Exec(n, b.networkId, b.address); err != nil {
		return false
	}

	return true
}

func (b *Block) SetBlockNumberWithoutDb(n int) {
	b.syncedBlockNumber = n
}

func New(networkId string, abiHash string, address string, syncedBlockNumber int) Block {
	return Block{
		networkId:         networkId,
		address:           abiHash,
		abiHash:           address,
		syncedBlockNumber: syncedBlockNumber,
	}
}

func (b *Block) ToJSON() map[string]interface{} {
	i := map[string]interface{}{}
	i["network_id"] = b.networkId
	i["address"] = b.address
	i["abi_hash"] = b.abiHash
	i["synced_block_number"] = b.syncedBlockNumber
	return i
}

func ParseJSON(blob map[string]interface{}) *Block {
	b := New(blob["network_id"].(string), blob["address"].(string), blob["abi_hash"].(string), int(blob["synced_block_number"].(float64)))
	return &b
}

func (b *Block) ToString() string {
	s := b.ToJSON()
	byt, err := json.Marshal(s)
	if err != nil {
		return ""
	}

	return string(byt)
}

func (b *Block) Save(socket *zmq.Socket) error {
	// Send hello.
	request := message.Request{
		Command: "categorizer_block_set",
		Param:   b.ToJSON(),
	}
	fmt.Println("Sending message to STATIC server to get abi. The mesage sent to server")
	fmt.Println(request.ToString())
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for abi getting from static controller")
		return fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller")
		return fmt.Errorf("receiving: %w", err)
	}

	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse abi reply")
		return fmt.Errorf("spaghetti block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure")
		return fmt.Errorf("spaghetti block reply status is not ok: %s", reply.Message)
	}

	fmt.Println("Abi build and returned")
	return nil
}

func GetAll(db *sql.DB, networkId string) []Block {
	var blocks []Block

	rows, err := db.Query("SELECT network_id, address, abi_hash, synced_block FROM categorizer_blocks WHERE network_id = ?", networkId)
	if err != nil {
		fmt.Println("Failed to query all categorizer blocks for network id ", networkId)
		fmt.Println(err.Error())
		return blocks
	}
	defer rows.Close()
	// An album slice to hold data from returned rows.

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var block Block
		if err := rows.Scan(&block.networkId, &block.address, &block.abiHash, &block.syncedBlockNumber); err != nil {
			rows.Close()
			return blocks
		}
		blocks = append(blocks, block)
	}
	if err = rows.Err(); err != nil {
		rows.Close()
		return blocks
	}

	rows.Close()

	return blocks
}

func RemoteGet(socket *zmq.Socket, networkId string, address string) (*Block, error) {
	// Send hello.
	request := message.Request{
		Command: "get",
		Param: map[string]interface{}{
			"networkId": networkId,
			"address":   address,
		},
	}
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		return nil, fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		return nil, fmt.Errorf("receiving from SDS Categorizer for 'get' command: %w", err)
	}

	fmt.Println(r)
	reply, err := message.ParseReply(r)
	if err != nil {
		return nil, fmt.Errorf("categorizer block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		return nil, fmt.Errorf("categorizer block reply status is not ok: %s", reply.Message)
	}

	returnedBlock := reply.Params["block"].(interface{})
	b := ParseJSON(returnedBlock.(map[string]interface{}))
	return b, nil
}

func RemoteGetAll(socket *zmq.Socket) ([]*Block, error) {
	// Send hello.
	request := message.Request{
		Command: "get_all",
		Param:   map[string]interface{}{},
	}
	if _, err := socket.SendMessage(request.ToString()); err != nil {
		fmt.Println("Failed to send a command for abi getting from static controller")
		return nil, fmt.Errorf("sending: %w", err)
	}

	// Wait for reply.
	r, err := socket.RecvMessage(0)
	if err != nil {
		fmt.Println("Failed to receive reply from static controller")
		return nil, fmt.Errorf("receiving: %w", err)
	}

	fmt.Println(r)
	reply, err := message.ParseReply(r)
	if err != nil {
		fmt.Println("Failed to parse abi reply")
		return nil, fmt.Errorf("spaghetti block invalid Reply: %w", err)
	}
	if !reply.IsOK() {
		fmt.Println("The static server returned failure")
		return nil, fmt.Errorf("spaghetti block reply status is not ok: %s", reply.Message)
	}

	returnedBlocks := reply.Params["blocks"].([]interface{})
	blocks := make([]*Block, len(returnedBlocks))
	for i, returnedBlock := range returnedBlocks {
		b := ParseJSON(returnedBlock.(map[string]interface{}))

		blocks[i] = b
	}

	return blocks, nil
}
