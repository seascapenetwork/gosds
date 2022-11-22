// the categorizer package keeps data types used by SDS Categorizer.
// the data type functions as well as method to obtain data from SDS Categorizer.
package categorizer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/blocklords/gosds/message"
	"github.com/blocklords/gosds/remote"
	"github.com/blocklords/gosds/static"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

////////////////////////////////////////////////////////////////////////////
//
// Abi struct is used for EVM based categorizer.
// it has the smartcontract interface to parse the raw spaghetti data into categorized data.
// its the wrapper over the SDS Static abi.
//
////////////////////////////////////////////////////////////////////////////
type Abi struct {
	staticAbi *static.Abi
	i         abi.ABI // interface
}

func (a *Abi) StringReader() *strings.Reader {
	return strings.NewReader(string(a.staticAbi.Bytes))
}

func (a *Abi) GetMethod(method string) (*abi.Method, error) {
	for _, m := range a.i.Methods {
		if m.Name == method {
			return &m, nil
		}
	}

	return nil, errors.New("method not found")
}

// given a transaction data, return a categorized variant
// first is returned a method name, second it returns method arguments.
func (a *Abi) ParseTxInput(data string) (string, map[string]interface{}, error) {
	inputs := map[string]interface{}{}

	offset := 0
	if len(data) > 2 && data[:2] == "0x" || data[:2] == "0X" {
		offset += 2
	}

	// decode txInput method signature
	decodedSig, err := hex.DecodeString(data[offset : 8+offset])
	if err != nil {
		return "", inputs, fmt.Errorf("failed to extract method signature from transaction data. the hex package error: %w", err)
	}

	// recover Method from signature and ABI
	method, err := a.i.MethodById(decodedSig)
	if err != nil {
		return "", inputs, fmt.Errorf("failed to find a method by its signature. geth package error: %w", err)
	}

	// decode txInput Payload
	decodedData, err := hex.DecodeString(data[8+offset:])
	if err != nil {
		return method.Name, inputs, fmt.Errorf("failed to extract method input arguments from transaction data. the hex package error: %w", err)
	}

	err = method.Inputs.UnpackIntoMap(inputs, decodedData)
	if err != nil {
		return method.Name, inputs, fmt.Errorf("failed to parse method input parameters into map. the geth package error: %w", err)
	}

	return method.Name, inputs, nil
}

// given abi hash and a json body, this method returns a new abi
func Build(abiHash string, body interface{}) (*Abi, error) {
	staticAbi := static.BuildAbi(body)
	staticAbi.AbiHash = abiHash

	abiObj := Abi{staticAbi: staticAbi}

	abiReader := abiObj.StringReader()
	i, err := abi.JSON(abiReader)
	if err != nil {
		return abiObj, fmt.Errorf("failed to parse body. probably an invalid json body. the geth package error: %w", err)
	}
	abiObj.i = i

	return &abiObj, nil
}

// returns an abi
func RemoteAbi(socket *remote.Socket, abiHash string) (*Abi, error) {
	// Send hello.
	request := message.Request{
		Command: "abi_get",
		Param: map[string]interface{}{
			"abi_hash": abiHash,
		},
	}

	params, err := socket.RequestRemoteService(&request)
	if err != nil {
		return nil, err
	}

	return Build(abiHash, params["abi"])
}
