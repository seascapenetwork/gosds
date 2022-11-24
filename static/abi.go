package static

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

type Abi struct {
	Bytes []byte
	// Body abi.ABI
	Body    interface{}
	AbiHash string
	exists  bool
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

func BuildAbi(body interface{}) *Abi {
	abi := Abi{Body: body, exists: false}
	byt, err := json.Marshal(body)
	if err != nil {
		return &abi
	}

	abi.Bytes = byt
	abi.CalculateAbiHash()

	return &abi
}

func BuildAbiFromBytes(bytes []byte) *Abi {
	body := []interface{}{}
	err := json.Unmarshal(bytes, &body)
	if err != nil {
		fmt.Println("Failed to convert abi bytes to body interface")
	}

	abi := Abi{Body: body, exists: false, Bytes: bytes}
	return &abi
}
