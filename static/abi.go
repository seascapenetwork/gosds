package static

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/ethereum/go-ethereum/crypto"
)

type Abi struct {
	Bytes []byte
	// Body abi.ABI
	Body    interface{}
	AbiHash string
	exists  bool
}

// Creates the JSON object with abi hash and abi body.
func (abi *Abi) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"abi":      abi.Body,
		"abi_hash": abi.AbiHash,
	}
}

// Creates the abi hash from the abi body
// The abi hash is the unique identifier of the abi
func (a *Abi) CalculateAbiHash() {
	hash := crypto.Keccak256Hash(a.Bytes)
	a.AbiHash = hash.String()[2:10]
}

// check whether the abi when its build was built from the database or in memory
func (a *Abi) Exists() bool {
	return a.exists
}

// set the exists flag
func (a *Abi) SetExists(exists bool) {
	a.exists = exists
}

// creates the Abi data based on the abi JSON. The function calculates the abi hash
// but won't set it in the database.
func NewAbi(body interface{}) (*Abi, error) {
	abi := Abi{Body: body, exists: false}
	byt, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	abi.Bytes = byt
	abi.CalculateAbiHash()

	return &abi, nil
}

// creates the Abi data based on the JSON string. This function calculates the abi hash
// but won't set it in the database.
func NewAbiFromBytes(bytes []byte) *Abi {
	body := []interface{}{}
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		fmt.Println("Failed to convert abi bytes to body interface")
	}

	abi := Abi{Body: body, exists: false, Bytes: bytes}
	return &abi
}

// Sends the ABI information to the remote SDS Static.
func RemoteAbiRegister(socket *remote.Socket, body interface{}) (map[string]interface{}, error) {
	// Send hello.
	request := message.Request{
		Command: "abi_register",
		Param: map[string]interface{}{
			"abi": body,
		},
	}

	return socket.RequestRemoteService(&request)
}

// Returns the abi from the remote server
func RemoteAbi(socket *remote.Socket, abi_hash string) (*Abi, error) {
	// Send hello.
	request := message.Request{
		Command: "abi_get",
		Param: map[string]interface{}{
			"abi_hash": abi_hash,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	abi_bytes, ok := params["abi"]
	if !ok {
		return nil, errors.New("missing 'abi' parameter from the SDS Static 'abi_get' command")
	}

	new_abi, err := NewAbi(abi_bytes)
	if err != nil {
		return nil, err
	}
	new_abi.AbiHash = abi_hash

	return new_abi, nil
}
