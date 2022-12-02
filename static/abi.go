package static

import (
	"encoding/json"
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

func (a *Abi) CalculateAbiHash() {
	hash := crypto.Keccak256Hash(a.Bytes)
	a.AbiHash = hash.String()[2:10]
}

func (a *Abi) Exists() bool {
	return a.exists
}

func (a *Abi) SetExists(exists bool) {
	a.exists = exists
}

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

	new_abi, err := NewAbi(params["abi"])
	if err != nil {
		return nil, err
	}
	new_abi.AbiHash = abi_hash

	return new_abi, nil
}
